package main

import "errors"

var ErrEmptyCommand = errors.New("пустая команда")
var ErrInvalidCommand = errors.New("неправильная команда")
