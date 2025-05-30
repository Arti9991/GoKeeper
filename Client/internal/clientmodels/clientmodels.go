package clientmodels

import "errors"

type CardInfo struct {
	Number  string
	ExpDate string
	CVVcode string
	Holder  string
}

type LoginInfo struct {
	Login    string
	Password string
}

type JournalInfo struct {
	Opperation string
	StorageID  string
	DataType   string
	MetaInfo   string
	SaveTime   string
	Sync       bool
}

type NewerData struct {
	StorageID string
	DataType  string
	MetaInfo  string
	SaveTime  string
	Data      []byte
}

var ErrorInput = errors.New("ошибка ввода")
var ErrBigData = errors.New("слишком большой размер данных")
var ErrBigCantSave = errors.New("не удалось сохранить файл")
var ErrNewerData = errors.New("более новые данные на сервере")
var ErrNoOfflineList = errors.New("cannot do this in offline")
var ErrNoSuchRows = errors.New("no rows for update")

var TokenFile = "./Token.txt"
var JournalFile = "./Journal.jl"
var StorageDir = "./Storage/"

var ErrBadLogin = errors.New("Слишком короткий логин")
var ErrBadPassowrd = errors.New("Пароль должен содержать не менее 6 символов")
