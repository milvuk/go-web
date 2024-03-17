package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
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
		log.Fatal("Invalid type assertion")
	}

	return value
}

func getDbHandle() *sql.DB {
	host := viperEnvVariable("DB_HOST")
	port := viperEnvVariable("DB_PORT")
	user := viperEnvVariable("DB_USER")
	pass := viperEnvVariable("DB_PASS")
	dbname := viperEnvVariable("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, dbname)

	var err error

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	log.Println("DB Connected!")
	return db
}

func main() {
	db := getDbHandle()

	listenAddr := viperEnvVariable("API_LISTEN_ADDR")

	srv := APIServer{
		listenAddr: listenAddr,
		db:         db,
	}

	srv.run()
}
