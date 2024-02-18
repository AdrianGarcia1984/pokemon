package env

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

func init (){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func EnvConfig(env, def string) string {
	value,err:=os.LookupEnv(env)
	if err{
		return value
	}
	return def	
}