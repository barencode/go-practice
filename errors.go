package main

import "errors"

var ErrMinerNotFound = errors.New("шахтёр не найден")
var ErrMineClosed = errors.New("шахта закрыта")
var ErrInvalidAmount = errors.New("неверное количество угля")
var ErrNotEnoughCoal = errors.New("недостаточно угля")
