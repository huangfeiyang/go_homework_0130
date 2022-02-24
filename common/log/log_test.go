package log

import (
	"errors"
	"testing"
)

func TestSetLevel(t *testing.T) {
	err := errors.New("this is error")
	debug := errors.New("this is debug")
	info := errors.New("this is info")

	SetLevel(0)
	Error(err)
	Debug(debug)
	Info(info)

	SetLevel(1)
	Error(err)
	Debug(debug)
	Info(info)

	SetLevel(2)
	Error(err)
	Debug(debug)
	Info(info)
}