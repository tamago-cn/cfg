package cfg

import (
	"testing"
)

type Demo struct {
	Host string `ini:"host" comment:"demo host"`
	Port int    `ini:"port" comment:"demo port"`
}
type Demo1 struct {
	Host string `ini:"host"`
	Port int    `ini:"port"`
}

func TestLoad(t *testing.T) {
	d := &Demo{
		Host: "127.0.0.1",
		Port: 6666,
	}
	d1 := &Demo1{
		Host: "127.0.0.1",
		Port: 8888,
	}
	RegistSection("demo", d, nil, nil)
	RegistSection("demo1", d1, nil, nil)
	err := Load("app.ini", true)
	if err != nil {
		t.Errorf("load config error: %s", err.Error())
	}
	err = Save()
	if err != nil {
		t.Errorf("load config error: %s", err.Error())
	}
}
