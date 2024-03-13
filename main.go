package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func viperEnvVariable(key string) string {
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	value, ok := viper.Get(key).(string)
	if !ok {
		log.Fatalf("Invalid type assertion")
	}

	return value
}

func main() {
	value := viperEnvVariable("DB_PORT")

	fmt.Println(value)
}
