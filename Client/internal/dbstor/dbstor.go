package dbstor

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Arti9991/GoKeeper/client/internal/clientmodels"
	_ "modernc.org/sqlite"
)

// структура с информацией о локальной базе данных
type DBStor struct {
	DbInfo string
	Db     *sql.DB
}

var (
	QuerryCreateTable = `CREATE TABLE IF NOT EXISTS datainfo (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		storage_id VARCHAR(64) NOT NULL UNIQUE,
		meta_info TEXT,
		data_type TEXT NOT NULL CHECK (data_type IN ('BINARY','CARD','TEXT', 'AUTH')),
		saved_time TIMESTAMP,
		sync BOOL
		);`
	QuerrySave = `INSERT INTO datainfo (storage_id, meta_info, data_type, saved_time, sync)
	VALUES  ($1, $2, $3, $4, $5);`
	QuerryGet = `SELECT meta_info, data_type, saved_time, sync FROM datainfo 
	WHERE storage_id = $1;`
	QuerryGetSync = `SELECT storage_id, meta_info, data_type, saved_time FROM datainfo 
	WHERE sync = FALSE;`
	QuerryDrop   = `DROP TABLE datainfo;`
	QuerryGetAll = `SELECT * FROM datainfo;`
	QuerryDone   = `UPDATE datainfo SET sync=TRUE
	WHERE storage_id = $1;`
	QuerryUndone = `UPDATE datainfo SET sync=FALSE
	WHERE storage_id = $1;`
	QuerryUpdate = `UPDATE datainfo SET meta_info=$2, data_type=$3, saved_time=$4, sync=$5
	WHERE storage_id = $1;`
	QuerryDelete = `DELETE FROM datainfo
  	WHERE storage_id = $1`
)

// DbInit инициализация локальной базы данных для хранения информации о данных
func DbInit(DBInfo string) (*DBStor, error) {
	var db DBStor
	var err error

	// открываем соединение
	db.Db, err = sql.Open("sqlite", DBInfo)
	if err != nil && DBInfo != "" {
		return &DBStor{}, err
	} else if DBInfo == "" {
		return &DBStor{}, errors.New("turning off data base mode by command dbinfo = _")
	}
	// проверяем соединение
	if err = db.Db.Ping(); err != nil {
		return &DBStor{}, err
	}
	// создаем таблицу
	_, err = db.Db.Exec(QuerryCreateTable)
	if err != nil {
		return &DBStor{}, err
	}

	return &db, nil
}

// ReinitTable функция очистки и пересоздания таблицы
// для функции logout (и тестов)
func (db *DBStor) ReinitTable() error {
	var err error

	_, err = db.Db.Exec(QuerryDrop)
	if err != nil {
		return err
	}

	_, err = db.Db.Exec(QuerryCreateTable)
	if err != nil {
		return err
	}

	return nil
}

// ShowTable функция вывода всех записей в таблице на экран
func (db *DBStor) ShowTable() error {
	var err error
	var Jr clientmodels.JournalInfo
	var id string
	var StorageID string
	var sended bool

	rows, err := db.Db.Query(QuerryGetAll)
	if err != nil {
		return err
	}
	defer rows.Close()
	fmt.Printf("%-10s %-64s %-25s %-10s %-40s %-6v\n", "ID", "StorageID", "MetaInfo", "Type", "SaveTime", "Sync")
	for rows.Next() {
		err = rows.Scan(&id, &StorageID, &Jr.MetaInfo, &Jr.DataType, &Jr.SaveTime, &sended)
		if err != nil {
			return err
		}
		fmt.Printf("%-10s %-64s %#-25v %-10s %-40s %#-6v\n", id, StorageID, Jr.MetaInfo, Jr.DataType, Jr.SaveTime, sended)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}

// SaveNew сохранение ифнормации о данных в таблицу
func (db *DBStor) SaveNew(StorageID string, Jr clientmodels.NewerData) error {

	var err error
	_, err = db.Db.Exec(QuerrySave, StorageID, Jr.MetaInfo, Jr.DataType, Jr.SaveTime, false)
	if err != nil {
		if strings.Contains(err.Error(), "2067") {
			err2 := db.UpdateInfoNewer(StorageID, Jr)
			if err2 != nil {
				return err2
			}
			return nil
		} else {
			return err
		}
	}
	return nil
}

// MarkDone отметка, что данные были отправлены на сервер
// и синхронизированны
func (db *DBStor) MarkDone(StorageID string) error {

	var err error
	_, err = db.Db.Exec(QuerryDone, StorageID)
	if err != nil {
		return err
	}
	return nil
}

// MarkUnDone отметка, что локальне данные нельзя считать синхронизированными
// (для случая когда метод Get возвращает более свежие данные с сервера)
func (db *DBStor) MarkUnDone(StorageID string) error {

	var err error
	_, err = db.Db.Exec(QuerryUndone, StorageID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateInfoNewer принудительное обновление данных в таблице с отметкой о синхронизации
func (db *DBStor) UpdateInfoNewer(StorageID string, Jr clientmodels.NewerData) error {

	var err error
	_, err = db.Db.Exec(QuerryUpdate, StorageID, Jr.MetaInfo, Jr.DataType, Jr.SaveTime, true)
	if err != nil {
		return err
	}
	return nil
}

// UpdateInfoNewer принудительное обновление данных в таблице
// без отметки о синхронизации
func (db *DBStor) UpdateInfo(StorageID string, Jr clientmodels.NewerData) error {

	res, err := db.Db.Exec(QuerryUpdate, StorageID, Jr.MetaInfo, Jr.DataType, Jr.SaveTime, false)
	if err != nil {
		if strings.Contains(err.Error(), "23505") {
			return err
		} else {
			return err
		}
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	} else if n == 0 {
		return clientmodels.ErrNoSuchRows
	}
	return nil
}

// Get получение значений из таблицы по StorageID
func (db *DBStor) Get(StorageID string) (clientmodels.NewerData, error) {

	var err error
	var out clientmodels.NewerData
	var Sync bool

	out.StorageID = StorageID

	row := db.Db.QueryRow(QuerryGet, StorageID)
	err = row.Scan(&out.MetaInfo, &out.DataType, &out.SaveTime, &Sync)
	if err != nil {
		return out, err
	}

	return out, nil
}

// GetForSync получение списка данных из таблицы, подлежащих синхронизацииэ
// (отметка о sync - false)
func (db *DBStor) GetForSync() ([]clientmodels.NewerData, error) {

	var err error
	var out []clientmodels.NewerData

	rows, err := db.Db.Query(QuerryGetSync)
	if err != nil {
		return out, err
	}
	defer rows.Close()

	for rows.Next() {
		var val clientmodels.NewerData
		err = rows.Scan(&val.StorageID, &val.MetaInfo, &val.DataType, &val.SaveTime)
		if err != nil {
			return out, err
		}
		out = append(out, val)
	}
	if err := rows.Err(); err != nil {
		return out, err
	}
	return out, nil
}

// DeleteData удаление строки из таблицы
func (db *DBStor) DeleteData(StorageID string) error {
	var err error
	_, err = db.Db.Exec(QuerryDelete, StorageID)
	if err != nil {
		return err
	}
	return nil
}
