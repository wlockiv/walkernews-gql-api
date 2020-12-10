package controllers

import "errors"

var NotFoundError = errors.New("the record was not found")
var WrongUsernameOrPasswordError = errors.New("the username or password in incorrect")
