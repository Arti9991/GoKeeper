package pgstor

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"github.com/Arti9991/GoKeeper/server/internal/logger"
	"github.com/Arti9991/GoKeeper/server/internal/server/servermodels"
)

var (
	QuerryCreateType = `CREATE TYPE types AS ENUM('BINARY','CARD','TEXT', 'AUTH');`
	QuerryCreateData = `CREATE TABLE IF NOT EXISTS datainfo (
		id SERIAL PRIMARY KEY,
		user_id VARCHAR(16) NOT NULL,
		storage_id VARCHAR(64) NOT NULL,
		meta_info TEXT,
		data_type types,
		saved_time TIMESTAMP,
		UNIQUE (user_id, storage_id)
		);`
	QuerrySaveDataInfo = `INSERT INTO datainfo
	(id, user_id, storage_id, meta_info, data_type, saved_time)
  	VALUES  (DEFAULT, $1, $2, $3, $4, $5);`
	QuerryGetDataTime = `SELECT saved_time FROM  datainfo
	WHERE storage_id = $1 AND user_id = $2`
	QuerryGetServDataInfo = `SELECT meta_info, data_type FROM  datainfo
	WHERE storage_id = $1 AND user_id = $2`
	QuerryUpdateDataInfo = `UPDATE datainfo
	SET meta_info = $3, data_type = $4, saved_time = $5
  	WHERE storage_id = $1 AND user_id = $2`
	QuerryGetDataInfo = `SELECT meta_info, data_type, saved_time FROM  datainfo
	WHERE storage_id = $1 AND user_id = $2`
	QuerryGetDataList = `SELECT storage_id, meta_info, data_type, saved_time FROM  datainfo
	WHERE user_id = $1`
	QuerryDeleteDataInfo = `DELETE FROM datainfo
  	WHERE storage_id = $1 AND user_id = $2`
)

// интерфеейс для базы данных с информацией о сохраненных данных
type InfoStorage interface {
	SaveNewData(userID string, DataInf servermodels.SaveDataInfo) (servermodels.SaveDataInfo, error)
	GetData(userID string, storageID string) (servermodels.SaveDataInfo, error)
	GetDataList(userID string) ([]servermodels.SaveDataInfo, error)
	UpdateData(userID string, DataInf servermodels.SaveDataInfo) error
	DeleteData(userID string, storageID string) error
	CloseDataDB() error
}

// DBStor структура для интерфейсов базы данных.
type DBStor struct {
	DB     *sql.DB // соединение с базой
	DBInfo string  // информация для подключения к базе
}

// DBinit инициализация хранилища и создание/подключение к таблице
// с информацией о данных.
func DBDataInit(DBInfo string) (*DBStor, error) {
	var db DBStor
	var err error

	db.DBInfo = DBInfo
	// открытие соединения с базой
	db.DB, err = sql.Open("pgx", DBInfo)
	if err != nil && DBInfo != "" {
		logger.Log.Error("Error in opening datainfo Db", zap.Error(err))
		return &DBStor{}, err
	} else if DBInfo == "" {
		return &DBStor{}, errors.New("turning off data base mode by command dbinfo = _")
	}
	// создание таблицы для информации о данных
	_, err = db.DB.Exec(QuerryCreateData)
	if err != nil {
		// если в таблице есть неопределенный тип, определяем его
		if strings.Contains(err.Error(), "SQLSTATE 42710") || strings.Contains(err.Error(), "SQLSTATE 42704") {
			// определеяем тип для хранения типа бинарных данных
			_, err = db.DB.Exec(QuerryCreateType)
			if err != nil {
				logger.Log.Error("Error in creating datainfo Type Db", zap.Error(err))
				return &DBStor{}, err
			}

			// создаем таблицу для информации о данных
			_, err = db.DB.Exec(QuerryCreateData)
			if err != nil {
				logger.Log.Error("Error in creating datainfo Db", zap.Error(err))
				return &DBStor{}, err
			}
			// провреям соединение
			if err = db.DB.Ping(); err != nil {
				logger.Log.Error("Error in ping datainfo Db", zap.Error(err))
				return &DBStor{}, err
			}
			logger.Log.Info("✓ connected to datainfo db! with new status type!")
			return &db, nil
		} else {
			logger.Log.Error("Error in creating datainfo Db", zap.Error(err))
			return &DBStor{}, err
		}
	}
	// провреям соединение
	if err = db.DB.Ping(); err != nil {
		logger.Log.Error("Error in ping datainfo Db", zap.Error(err))
		return &DBStor{}, err
	}

	logger.Log.Info("✓ connected to datainfo db!")
	return &db, nil
}

// SaveNewData функция сохранения данных. Если в базе есть более свежие (по времени сохранения) данные,
// то возвращает их в ответ с соответствующим флагом
func (db *DBStor) SaveNewData(userID string, DataInf servermodels.SaveDataInfo) (servermodels.SaveDataInfo, error) {
	var err error
	var outData servermodels.SaveDataInfo
	// выполняем вставку полученной информации
	_, err = db.DB.Exec(QuerrySaveDataInfo, userID, DataInf.StorageID, DataInf.MetaInfo, DataInf.Type, DataInf.SaveTime)
	if err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			// если для данного пользователя и storageID уже есть запись
			// то получем время сохранения, указанное в этой записи
			var BaseTime time.Time
			row := db.DB.QueryRow(QuerryGetDataTime, DataInf.StorageID, userID)
			err = row.Scan(&BaseTime)
			if err != nil {
				logger.Log.Error("Error in getting saved time form datainfo Db", zap.Error(err))
				return outData, err
			}
			if BaseTime.After(DataInf.SaveTime) {
				// если время сохранения в записи из базы свежее, чем в запросе,
				// получаем исходную запись из базы и возвращаем ее
				row := db.DB.QueryRow(QuerryGetServDataInfo, DataInf.StorageID, userID)
				err = row.Scan(&outData.MetaInfo, &outData.Type)
				if err != nil {
					logger.Log.Error("Error in getting saved time form datainfo Db", zap.Error(err))
					return outData, err
				}
				outData.SaveTime = BaseTime
				outData.UserID = userID
				outData.StorageID = DataInf.StorageID
				return outData, servermodels.ErrNewerData
			} else {
				// если же полученное в запросе время свежее, чем в базе, то сохраняем информацию из запроса
				_, err = db.DB.Exec(QuerryUpdateDataInfo, DataInf.StorageID, userID, DataInf.MetaInfo, DataInf.Type, DataInf.SaveTime)
				if err != nil {
					logger.Log.Error("Error in update newer data to datainfo Db", zap.Error(err))
					return outData, err
				}
				return outData, nil
			}
		} else {
			logger.Log.Error("Error in saving data to datainfo Db", zap.Error(err))
			return outData, err
		}
	}
	return outData, nil
}

// GetData функция получения информации о конкретных данных для пользователя
func (db *DBStor) GetData(userID string, storageID string) (servermodels.SaveDataInfo, error) {
	var err error
	var outData servermodels.SaveDataInfo

	row := db.DB.QueryRow(QuerryGetDataInfo, storageID, userID)
	err = row.Scan(&outData.MetaInfo, &outData.Type, &outData.SaveTime)
	if err != nil {
		logger.Log.Error("Error in get data from datainfo Db", zap.Error(err))
		return outData, err
	}

	return outData, nil
}

// GetDataList функция получения списка для всех данных у пользователя
func (db *DBStor) GetDataList(userID string) ([]servermodels.SaveDataInfo, error) {
	var err error
	var outData []servermodels.SaveDataInfo

	rows, err := db.DB.Query(QuerryGetDataList, userID)
	if err != nil {
		logger.Log.Error("Error in get rows from datainfo Db", zap.Error(err))
		return outData, err
	}
	defer rows.Close()

	for rows.Next() {
		var dataLine servermodels.SaveDataInfo
		err = rows.Scan(&dataLine.StorageID, &dataLine.MetaInfo, &dataLine.Type, &dataLine.SaveTime)
		if err != nil {
			logger.Log.Error("Error in get data from datainfo Db", zap.Error(err))
			return outData, err
		}
		outData = append(outData, dataLine)
	}

	return outData, nil
}

// UpdateData фнукция для принудительного обновления данных, не взирая на время сохранения
func (db *DBStor) UpdateData(userID string, DataInf servermodels.SaveDataInfo) error {
	var err error

	_, err = db.DB.Exec(QuerryUpdateDataInfo, DataInf.StorageID, userID, DataInf.MetaInfo, DataInf.Type, DataInf.SaveTime)
	if err != nil {
		logger.Log.Error("Error in update newer data to datainfo Db", zap.Error(err))
		return err
	}
	return nil
}

// DeleteData функция для полного удаления информации о данных из базы
func (db *DBStor) DeleteData(userID string, storageID string) error {
	var err error

	_, err = db.DB.Exec(QuerryDeleteDataInfo, storageID, userID)
	if err != nil {
		logger.Log.Error("Error in delete data from datainfo Db", zap.Error(err))
		return err
	}
	return nil
}

// CloseDataDB функция закрытия соединения с базой информации о данных
func (db *DBStor) CloseDataDB() error {
	return db.DB.Close()
}
