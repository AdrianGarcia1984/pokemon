package database

import (
	"log"
	"os"
	"batalla_pokemon/env"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var host = env.EnvConfig("POSTGRES_HOST",os.Getenv("POSTGRES_HOST"))
//var host = os.Getenv("POSTGRES_HOST")

var user = os.Getenv("POSTGRES_USER")
var password = os.Getenv("POSTGRES_PASSWORD")
var dbname = os.Getenv("POSTGRES_DATABASE")
var port = os.Getenv("PORT")

var DSNENV = ("host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port)

var DB *gorm.DB

func DBConection() {

	var err error
	DB, err = gorm.Open(postgres.Open(DSNENV), &gorm.Config{})

	if err != nil {
		log.Println("error conexion database")
		log.Fatal(err)
		
	}
	log.Println("db connected")
}

