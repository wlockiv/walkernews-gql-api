package errors

import "errors"

func NewBaseError(err error) error {
	msg := "Error :: " + err.Error()
	return errors.New(msg)
}

func NewAuthError(err error) error {
	newErr := errors.New("Authentication - " + err.Error())
	return NewBaseError(newErr)
}

func NewDBError(query string, err error) error {
	newErr := errors.New("DB Query (" + query + ") - " + err.Error())
	return NewBaseError(newErr)
}

func NewUnmarshallError(item string, err error) error {
	newErr := errors.New("Unmarshalling Error (" + item + ") - " + err.Error())
	return NewBaseError(newErr)
}
