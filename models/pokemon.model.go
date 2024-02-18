package models

type Pokemon struct{
	Id int `json:"id"`
	Name string `json:"name"`
	Hp int `json:"hp"`
	Speed int `json:"speed"`
	Attack int `json:"attack"`
	Defense int `json:"defense"`
}


type BatleTable struct{
	Id int `json:"id"`
	Pokemon1 int `json:"pokemon_1"`
	Pokemon2 int`json:"pokemon_2"`
	ResultBatle int `json:"result_batle"`
}

