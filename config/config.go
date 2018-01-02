package config

import (
	"log"

	ini "gopkg.in/ini.v1"
)

var Config struct {
	Server    Server           `ini:"server"`
	DB        DBConfig         `ini:"db"`
	Grades    Grades           `ini:"grades"`
	Auth      Authentification `ini:"auth"`
	Mail      Mail             `ini:"mail"`
	Sessioner Sessioner        `ini:"sessioner"`
}

type Sessioner struct {
	StoragePath string `ini:"storagepath"`
}

type Server struct {
	Address string `ini:"address"`
	Ip      string `ini:"ip"`
	Port    int    `ini:"port"`
	DevMode bool   `ini:"devmode"`
}

type Mail struct {
	Sender   string `ini:"sender"`
	Address  string `ini:"address"`
	Port     int    `ini:"port"`
	Username string `ini:"username"`
	Password string `ini:"password"`
}

type Authentification struct {
	MinPwLength   uint   `ini:"minpwlength"`
	MaxPwLength   uint   `ini:"maxpwlength"`
	AdminPassword string `ini:"adminpassword"`
}

type Grades struct {
	Min uint `ini:"min"`
	Max uint `ini:"max"`
}

type DBConfig struct {
	Host     string `ini:"host"`
	Port     uint   `ini:"port"`
	User     string `ini:"user"`
	Password string `ini:"password"`
	Name     string `ini:"name"`
}

func Load(path string) {
	f, err := ini.Load(path)
	if err != nil {
		log.Fatal(err)
	}
	err = f.MapTo(&Config)
	if err != nil {
		log.Fatal(err)
	}
}
