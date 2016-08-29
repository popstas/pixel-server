package main

import (
	"os"
	"strconv"
	"log"
)

type Config struct {
	SerialPort  string
	SerialSpeed int
	WebHost     string
	WebPort     int
}

func setIntEnvvar(v *int, envName string){
	envValue := os.Getenv(envName)
	if envValue != ""{
		if envIntValue, err := strconv.Atoi(envValue); err != nil{
			log.Fatalf("Cannot convert value of envvar %s to int: %s", envName, envValue)
		} else {
			*v = envIntValue
		}
	}
}

func setStringEnvvar(v *string, envName string){
	envValue := os.Getenv(envName)
	if envValue != ""{
		*v = envValue
	}
}

