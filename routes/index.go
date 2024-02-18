package routes

import (
	"batalla_pokemon/database"
	"batalla_pokemon/models"
	"encoding/json"
	"fmt"

	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mtslzr/pokeapi-go"
)

type PokemonWin struct {
	Id           int `json:"id"`
	Result_batle int `json:"Result_batle"`
}

func PostPokemonHandler(w http.ResponseWriter, r *http.Request) {
	var pokemon, newPokemon models.Pokemon
	json.NewDecoder(r.Body).Decode(&pokemon)
	database.DB.Where(&models.Pokemon{Name: pokemon.Name}).Find(&newPokemon)
	if newPokemon.Id != 0 {
		w.WriteHeader(http.StatusMethodNotAllowed) //405
		w.Write([]byte("no se puede crear el pokemon, ya hay uno registrado con los datos ingresados"))
		return
	}
	createdPokemon := database.DB.Save(&pokemon)
	err := createdPokemon.Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) //400
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(&pokemon)
	w.WriteHeader(http.StatusOK) //201
	w.Write([]byte("pokemon creado con exito"))

}

func BatallaPokemonHandler(w http.ResponseWriter, r *http.Request) {
	var pokemon1, pokemon2 models.Pokemon
	var batleTable models.BatleTable
	params := mux.Vars(r)
	database.DB.First(&pokemon1, params["pokemon1"])
	database.DB.First(&pokemon2, params["pokemon2"])

	if pokemon1.Id == 0 || pokemon2.Id == 0 {
		w.WriteHeader(http.StatusBadRequest) //400
		w.Write([]byte("un pokemon no esta generado por favor verificar"))
		return
	}
	var pokemonWin = pokemonBatle(pokemon1.Id, pokemon2.Id)
	if pokemonWin == 0 {
		pokemonWin = startPokemonBatle(pokemon1.Id, pokemon2.Id)
		batleTable = models.BatleTable{
			Pokemon1: pokemon1.Id,
			Pokemon2: pokemon2.Id,
			ResultBatle: pokemonWin,
		}
		createdPokemon := database.DB.Save(&batleTable)
		err := createdPokemon.Error
		if err != nil {
			w.WriteHeader(http.StatusBadRequest) //400
			w.Write([]byte(err.Error()))
			return
		}
		json.NewEncoder(w).Encode(&batleTable)
		w.WriteHeader(http.StatusOK) //201
		w.Write([]byte("batalla generada con exito"))
	}
	w.WriteHeader(http.StatusOK) //200
	w.Write([]byte("batalla retornada con exito, el ganador es: "))
	w.Write([]byte(strconv.Itoa(pokemonWin)))
}

func pokemonBatle(pokemon1 int, pokemon2 int) int {
	if pokemon1 == 0 || pokemon2 == 0 {
		return 0
	}
	pokemonWin := PokemonWin{}
	database.DB.Raw("select Result_batle, id from batle_tables WHERE batle_tables.Pokemon1=" + strconv.Itoa(pokemon1) + " and batle_tables.Pokemon2=" + strconv.Itoa(pokemon2)).Scan(&pokemonWin)
	if pokemonWin.Id == 0 {
		return pokemonWin.Id
	}
	fmt.Println("desde PokemonBatle ", pokemonWin.Result_batle)
	return pokemonWin.Result_batle
}

func startPokemonBatle(p1 int, p2 int) int {
	if p1 == 0 || p2 == 0 {
		return 0
	}
	var pokemon1, pokemon2 models.Pokemon
	database.DB.First(&pokemon1, p1)
	database.DB.First(&pokemon2, p2)
	l1, err1 := pokeapi.Pokemon(pokemon1.Name)
	if err1 != nil {
		fmt.Println("err: ", err1.Error())
	}
	l2, err2 := pokeapi.Pokemon(pokemon2.Name)
	if err2 != nil {
		fmt.Println("err: ", err2.Error())
	}
	if pokemon1.Name == l1.Name {
		pokemon1.Attack = l1.Stats[1].BaseStat
		pokemon1.Hp = l1.Stats[0].BaseStat
		pokemon1.Defense = l1.Stats[2].BaseStat
		pokemon1.Speed = l1.Stats[5].BaseStat
		pokemon2.Attack = l2.Stats[1].BaseStat
		pokemon2.Hp = l2.Stats[0].BaseStat
		pokemon2.Defense = l2.Stats[2].BaseStat
		pokemon2.Speed = l2.Stats[5].BaseStat
	}
	if pokemon2.Name == l1.Name {
		pokemon2.Attack = l1.Stats[1].BaseStat
		pokemon2.Hp = l1.Stats[0].BaseStat
		pokemon2.Defense = l1.Stats[2].BaseStat
		pokemon2.Speed = l1.Stats[5].BaseStat
		pokemon1.Attack = l2.Stats[1].BaseStat
		pokemon1.Hp = l2.Stats[0].BaseStat
		pokemon1.Defense = l2.Stats[2].BaseStat
		pokemon1.Speed = l2.Stats[5].BaseStat
	}
	var pokemon1Hp = pokemon1.Hp
	var pokemon2Hp = pokemon2.Hp
	if pokemon1.Speed > pokemon2.Speed {
		var turn =1
		for pokemon1Hp > 0 && pokemon2Hp > 0 {	
			switch {
				case turn==1:
					pokemon1Hp = pokemon1Hp - (pokemon2.Attack - pokemon1.Defense)
					turn =2
				
				case turn==2:
					pokemon2Hp = pokemon2Hp - (pokemon1.Attack - pokemon2.Defense)
					turn =1					
			}
		}	
	}else{
		var turn =2
		for pokemon1Hp > 0 && pokemon2Hp > 0 {			
			switch {
				case turn==1:
					pokemon1Hp = pokemon1Hp - (pokemon2.Attack - pokemon1.Defense)
					turn =2
				
				case turn==2:
					pokemon2Hp = pokemon2Hp - (pokemon1.Attack - pokemon2.Defense)
					turn =1					
			}
		}
	}
	if pokemon1Hp <=0{
		return pokemon2.Id
	}else{
		return pokemon1.Id
	}
}
