package config

import (
	"github.com/pelletier/go-toml/v2"
	"io/ioutil"
)

type Config struct {
	Env    string       `toml:"env"`
	Goploy GoployConfig `toml:"goploy"`
	Log    LogConfig    `toml:"log"`
	Web    WebConfig    `toml:"web"`
}

type LogConfig struct {
	Path  string `toml:"path"`
	Split bool   `toml:"split"`
}

type WebConfig struct {
	Port string `toml:"port"`
}

type GoployConfig struct {
	ReportURL string `toml:"reportURL"`
	Key       string `toml:"key"`
	UIDType   string `toml:"uidType"`
	UID       string `toml:"uid"`
}

var Toml Config

func Create(filename string) {
	config, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = toml.Unmarshal(config, &Toml)
	if err != nil {
		panic(err)
	}
}

func Write(filename string, cfg Config) error {
	yamlData, err := toml.Marshal(&cfg)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, yamlData, 0644)
	if err != nil {
		return err
	}
	return nil
}
