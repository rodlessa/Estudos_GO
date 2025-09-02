package main

import (
	"fmt"
)

type Cliente struct {
	ID            int
	Nome          string
	Identificador int
	Address       string
}

func ListClientes() ([]Cliente, error) {
	db := connectdb() // pega a conexão já criada

	rows, err := db.Query("SELECT id, nome, identificador, address FROM clientes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clientes []Cliente
	for rows.Next() {
		var c Cliente
		if err := rows.Scan(&c.ID, &c.Nome, &c.Identificador, &c.Address); err != nil {
			return nil, err
		}
		clientes = append(clientes, c)
	}

	return clientes, nil
}

func PrintClientes() {
	clientes, err := ListClientes()
	if err != nil {
		fmt.Println("Erro ao listar clientes:", err)
		return
	}

	for _, c := range clientes {
		fmt.Printf("%d - %s - %d - %s\n", c.ID, c.Nome, c.Identificador, c.Address)
	}
}
