package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
)

type placamae struct {
	Marcas         []string `json:"marcas"`
	Modelos_mother []string `json:"modelos_mother"`
}
type GPUS struct {
	Marcas      []string            `json:"marcas"`
	Modelos_GPU map[string][]string `json:"modelos_gpu"`
}
type CPUS struct {
	Modelos_CPU map[string][]string `json:"modelos_cpu"`
}
type memorias struct {
	Memorias_tipo []string `json:"memoria_tipo"`
	Memorias      []string `json:"memoria"`
}
type ssd struct {
	SSD []string `json:"tamanho_ssd"`
}

func randomizer_mother() (string, string) {
	jsonFile, err := os.Open("produtos.json")
	if err != nil {
		log.Fatal(err)
	}
	var data placamae
	if err := json.NewDecoder(jsonFile).Decode(&data); err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	return data.Marcas[rand.Intn(len((data.Marcas)))], data.Modelos_mother[rand.Intn(len(data.Modelos_mother))]

}
func randomizer_gpus() (string, string, string) {
	jsonFile, err := os.Open("produtos.json")
	if err != nil {
		log.Fatal(err)
	}
	var data GPUS
	if err := json.NewDecoder(jsonFile).Decode(&data); err != nil {
		log.Fatal(err)
	}

	keys := make([]string, 0, len(data.Modelos_GPU))
	for k := range data.Modelos_GPU {
		keys = append(keys, k)
	}
	randKey := keys[rand.Intn(len(keys))]
	modelos := data.Modelos_GPU[randKey]
	modeloEscolhido := modelos[rand.Intn(len(modelos))]
	defer jsonFile.Close()
	return data.Marcas[rand.Intn(len((data.Marcas)))], randKey, modeloEscolhido

}
func randomizer_cpus() (string, string) {
	jsonFile, err := os.Open("produtos.json")
	if err != nil {
		log.Fatal(err)
	}
	var data CPUS
	if err := json.NewDecoder(jsonFile).Decode(&data); err != nil {
		log.Fatal(err)
	}

	keys := make([]string, 0, len(data.Modelos_CPU))
	for k := range data.Modelos_CPU {
		keys = append(keys, k)
	}
	randKey := keys[rand.Intn(len(keys))]
	modelos := data.Modelos_CPU[randKey]
	modeloEscolhido := modelos[rand.Intn(len(modelos))]
	defer jsonFile.Close()
	return randKey, modeloEscolhido

}
func randomizer_memoria() (string, string) {
	jsonFile, err := os.Open("produtos.json")
	if err != nil {
		log.Fatal(err)
	}
	var data memorias
	if err := json.NewDecoder(jsonFile).Decode(&data); err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	return data.Memorias_tipo[rand.Intn(len(data.Memorias_tipo))], data.Memorias[rand.Intn(len(data.Memorias))]

}
func randomizer_ssd() string {
	jsonFile, err := os.Open("produtos.json")
	if err != nil {
		log.Fatal(err)
	}
	var data ssd
	if err := json.NewDecoder(jsonFile).Decode(&data); err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	return data.SSD[rand.Intn(len(data.SSD))]

}
func insert_produtos_cpus() {
	for {
		connectdb()
		marc, modelo := randomizer_cpus()
		cpus := fmt.Sprintf("%s %s", marc, modelo)
		query :=
			`
		SELECT id, nome, marca, qnt, valor::float8 FROM produtos WHERE nome = $1
	`
		row := db.QueryRow(query, cpus)
		var id int
		var nome, marca string
		var qnt int
		var valor float64

		err := row.Scan(&id, &nome, &marca, &qnt, &valor)
		if err != nil {
			if err == sql.ErrNoRows {
				nome = cpus
				qnt := rand.Intn(20)
				valor := rand.Float64() * 5000.00
				fmt.Printf("Nenhum produto encontrado...\n Inserindo produto: %s Valor: R$ %.2f\n", cpus, valor)

				subQuery :=
					`
			insert into produtos (nome, marca, qnt, valor)
			values ($1,$2,$3,$4)
			`
				_, err := db.Exec(subQuery, nome, marc, qnt, valor)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("Produto inserido!")
			} else {
				log.Fatal(err)
			}
		} else {

			if qnt == 0 {
				novaQnt := rand.Intn(20) + 1
				subQuery :=
					`
			update produtos set qnt = $1 where nome = $2
			`
				_, err := db.Exec(subQuery, novaQnt, nome)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("Produto inserido!")
			} else {
				fmt.Println("Produtos encontrado com estoque, não há mais o que inserir")
				continue

			}

		}
	}

}

func main() {
	insert_produtos_cpus()
}
