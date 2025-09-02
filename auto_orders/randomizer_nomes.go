package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
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

func insert_nomes(db *pgx.Conn) {
	start := time.Now()
	totalInserts := 10000

	rows := make([][]interface{}, totalInserts)
	for i := 0; i < totalInserts; i++ {
		nome := gerar_nomes()
		enderecos := gerar_end()
		identificador := gerar_indet()
		rows[i] = []interface{}{nome, identificador, enderecos}
	}
	copyCount, err := db.CopyFrom(
		context.Background(),
		pgx.Identifier{"clientes"},
		[]string{"nome", "identificador", "address"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		log.Fatal("Erro no CopyFrom:", err)
	}

	elapsed := time.Since(start)
	insertsPerSec := float64(copyCount) / elapsed.Seconds()

	fmt.Printf("\n===== Estatísticas =====\n")
	fmt.Printf("Tempo total: %s\n", elapsed)
	fmt.Printf("Inserts por segundo: %.2f\n", insertsPerSec)
	fmt.Printf("========================\n")
}

func main() {
	conn := connectDB()                    // retorna *pgx.Conn
	defer conn.Close(context.Background()) // fecha no final

	insert_nomes(conn)
}
