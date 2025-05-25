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

var ErrorInput = errors.New("ошибка ввода")
