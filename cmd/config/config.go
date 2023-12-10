package config

import (
	"flag"
)

var Config struct {
	Port        string
	BaseAddress string
}

func ConfigPort() {
	flag.StringVar(&Config.Port, "a", ":8080", "Addres local port")

}

func ConfigBaseAddress() {
	flag.StringVar(&Config.BaseAddress, "b", "http://local", "Base local addres")

}
