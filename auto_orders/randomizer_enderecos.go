package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"strconv"
)

type enderecos struct {
	Enderecos []string `json:"nome"`
	Tipo_end  []string `json:"tipo"`
}

func gerar_end() string {
	jsonFile, err := os.Open("enderecos.json")
	if err != nil {
		log.Fatal(err)
	}
	var data enderecos
	if err := json.NewDecoder(jsonFile).Decode(&data); err != nil {
		log.Fatal(err)
	}
	tipo := data.Tipo_end[rand.Intn(len(data.Tipo_end))]
	rua1 := data.Enderecos[rand.Intn(len(data.Enderecos))]
	rua2 := data.Enderecos[rand.Intn(len(data.Enderecos))]
	num := strconv.Itoa(rand.Intn(500) + 1)
	return tipo + " " + rua1 + " " + rua2 + ", " + num
}
