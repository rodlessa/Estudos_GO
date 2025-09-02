package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

type nomes struct {
	Nomes      []string `json:"nome"`
	Sobrenomes []string `json:"sobrenome"`
}

func gerar_nomes() string {
	jsonFile, err := os.Open("nomes_sobrenomes.json")
	if err != nil {
		log.Fatal(err)
	}
	var data nomes
	if err := json.NewDecoder(jsonFile).Decode(&data); err != nil {
		log.Fatal(err)
	}

	nome_final := fmt.Sprintf("%s %s %s", data.Nomes[rand.Intn(len(data.Nomes))], data.Sobrenomes[rand.Intn(len(data.Sobrenomes))], data.Sobrenomes[rand.Intn(len(data.Sobrenomes))])
	return nome_final
}

func gerar_indet() string {
	cpf := []int{}
	for i := 0; i < 11; i++ {
		cpf = append(cpf, rand.Intn(10))
	}
	s := ""
	for _, n := range cpf {
		s += fmt.Sprintf("%d", n)
	}
	return s
}

func insert_nomes() {
	start := time.Now()
	totalInserts := 10000
	for i := 0; i < totalInserts; i++ {
		nome := gerar_nomes()
		identificador := gerar_indet()
		endereco := gerar_end()

		var id int
		err := db.QueryRow(`
        INSERT INTO clientes (nome, identificador, address)
        VALUES ($1, $2, $3)
        RETURNING id
    `, nome, identificador, endereco).Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
	}
	elapsed := time.Since(start)
	insertsPerSec := float64(totalInserts) / elapsed.Seconds()
	fmt.Printf("\n===== Estatísticas =====\n")
	fmt.Printf("Tempo total: %s\n", elapsed)
	fmt.Printf("Inserts por segundo: %.2f\n", insertsPerSec)
	fmt.Printf("========================\n")
}

func main() {
	connectdb()
	insert_nomes()
}
