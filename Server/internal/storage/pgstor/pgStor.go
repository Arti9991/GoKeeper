package pgstor

// import (
// 	"database/sql"
// 	"errors"

// 	"github.com/Arti9991/GoKeeper/server/internal/logger"
// )

// var (
// 	QuerryCreate = `CREATE TABLE IF NOT EXISTS keeper (
// 		id SERIAL PRIMARY KEY,
// 		user_id VARCHAR(16),
// 		user_login 	VARCHAR(100) NOT NULL UNIQUE,
// 		user_password VARCHAR(64) NOT NULL,
// 		storage_id VARCHAR(32) NOT NULL UNIQUE,
// 		meta_info TEXT,
// 		saved_time TIMESTAMP,
// 		registr_data_flag BOOLEAN NOT NULL DEFAULT FALSE,
// 		card_info_flag BOOLEAN NOT NULL DEFAULT FALSE,
// 		text_info_flag BOOLEAN NOT NULL DEFAULT FALSE,
// 		binary_data_flag BOOLEAN NOT NULL DEFAULT FALSE,
// 		);`
// 	QuerryNewUser = `INSERT INTO keeper
// 	(id, user_id, user_login, user_password, storage_id, saved_time)
//   	VALUES  (DEFAULT, $1, $2, $3, $4, $5);`
// )

// // DBStor структура для интерфейсов базы данных.
// type DBStor struct {
// 	DB     *sql.DB // соединение с базой
// 	DBInfo string  // информация для подключения к базе
// }

// // DBinit инициализация хранилища и создание/подключение к таблице.
// func DBinit(DBInfo string) (*DBStor, error) {
// 	var db DBStor
// 	var err error

// 	db.DBInfo = DBInfo

// 	db.DB, err = sql.Open("pgx", DBInfo)
// 	if err != nil && DBInfo != "" {
// 		return &DBStor{}, err
// 	} else if DBInfo == "" {
// 		return &DBStor{}, errors.New("turning off data base mode by command dbinfo = _")
// 	}

// 	if err = db.DB.Ping(); err != nil {
// 		return &DBStor{}, err
// 	}

// 	_, err = db.DB.Exec(QuerryCreate)
// 	if err != nil {
// 		return &DBStor{}, err
// 	}
// 	logger.Log.Info("✓ connected to ShortURL db!")
// 	return &db, nil
// }

// func (db *DBStor) SaveNewUser(userID string, userLogin string, userPassw string) error {
// 	var err error

// 	_, err = db.DB.Exec(QuerryNewUser, userID, userLogin, userPassw, "registration", )
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
