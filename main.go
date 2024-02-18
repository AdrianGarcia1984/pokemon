package main

import (
	//"fmt"
	"batalla_pokemon/database"
	"batalla_pokemon/models"
	"batalla_pokemon/routes"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/rs/cors"
)

func main() {

	database.DBConection()
	database.DB.AutoMigrate(models.Pokemon{}, models.BatleTable{})

	r := mux.NewRouter()

	r.HandleFunc("/pokemon-ba+le", routes.PostPokemonHandler).Methods("POST")
	r.HandleFunc("/pokemon-ba+le/batle/{pokemon1}/{pokemon2}", routes.BatallaPokemonHandler).Methods("GET")

	handler := cors.Default().Handler(r)

	http.ListenAndServe(":3001", handler)
}
