package log

import "testing"

func TestInitLogger(t *testing.T) {
	InitLogger()
	Info("this is info message,username is %s","watertreestar")
}