package pgstor

import (
	"database/sql"
	"errors"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Arti9991/GoKeeper/server/internal/logger"
	"github.com/Arti9991/GoKeeper/server/internal/server/servermodels"
)

var (
	QuerryCreateUsers = `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		user_id VARCHAR(16),
		user_login 	VARCHAR(100) NOT NULL UNIQUE,
		user_password VARCHAR(64) NOT NULL
		);`
	QuerryNewUser = `INSERT INTO users (id, user_id, user_login, user_password)
  	VALUES  (DEFAULT, $1, $2, $3);`
	QuerryGetUser = `SELECT user_password FROM users 
	WHERE user_login = $1;`
)

// DBUsersStor структура для.
type DBUsersStor struct {
	DB     *sql.DB // соединение с базой
	DBInfo string  // информация для подключения к базе
}

// DBinit инициализация хранилища и создание/подключение к таблице.
func DBUsersinit(DBInfo string) (*DBUsersStor, error) {
	var db DBUsersStor
	var err error

	db.DBInfo = DBInfo

	db.DB, err = sql.Open("pgx", DBInfo)
	if err != nil && DBInfo != "" {
		return &DBUsersStor{}, err
	} else if DBInfo == "" {
		return &DBUsersStor{}, errors.New("turning off data base mode by command dbinfo = _")
	}

	if err = db.DB.Ping(); err != nil {
		return &DBUsersStor{}, err
	}

	_, err = db.DB.Exec(QuerryCreateUsers)
	if err != nil {
		return &DBUsersStor{}, err
	}
	logger.Log.Info("✓ connected to Users db!")
	return &db, nil
}

func (db *DBUsersStor) SaveNewUser(userID string, userLogin string, userPassw string) error {
	var err error

	_, err = db.DB.Exec(QuerryNewUser, userID, userLogin, userPassw)
	if err != nil {
		return err
	}
	return nil
}

func (db *DBUsersStor) GetUser(userLogin string) (string, error) {
	var err error
	var pass string

	row := db.DB.QueryRow(QuerryGetUser, userLogin)
	err = row.Scan(&pass)
	if err != nil && strings.Contains(err.Error(), "no rows") {
		return "", servermodels.ErrorNoSuchUser
	} else if err != nil {
		return "", err
	}

	return pass, nil
}
