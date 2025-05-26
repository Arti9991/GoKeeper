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
}

var ErrorInput = errors.New("ошибка ввода")

var TokenFile = "./Token.txt"
var JournalFile = "./Journal.txt"
