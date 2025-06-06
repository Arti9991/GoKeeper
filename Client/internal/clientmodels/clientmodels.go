package clientmodels

import "errors"

// структура с информацией о карте
type CardInfo struct {
	Number  string
	ExpDate string
	CVVcode string
	Holder  string
}

// структура с информацией авторизации
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

// структура с данными и информацией о них
type NewerData struct {
	StorageID string
	DataType  string
	MetaInfo  string
	SaveTime  string
	Data      []byte
}

// служебные ошибки
var (
	ErrBigCantSave   = errors.New("не удалось сохранить файл")
	ErrNewerData     = errors.New("более новые данные на сервере")
	ErrNoOfflineList = errors.New("cannot do this in offline")
	ErrNoSuchRows    = errors.New("no rows for update")
)

// пути для вспомогательных файлов
var TokenFile = "./Token.txt"
var JournalFile = "./Journal.jl"
var StorageDir = "./Storage/"

// ошибки ввода
var (
	ErrorInput     = errors.New("ошибка ввода")
	ErrBigData     = errors.New("слишком большой размер данных")
	ErrBadLogin    = errors.New("Слишком короткий логин")
	ErrBadPassowrd = errors.New("Пароль должен содержать не менее 6 символов")
	ErrShortNum    = errors.New("слишком короткий номер карты")
	ErrBadCVV      = errors.New("cvv должен содержать три символа")
	ErrUserAbort   = errors.New("отмена операции пользователем")
)
