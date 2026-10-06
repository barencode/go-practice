package main

import "errors"

var ErrMinerNotFound = errors.New("шахтёр не найден")
var ErrMineClosed = errors.New("шахта закрыта")
