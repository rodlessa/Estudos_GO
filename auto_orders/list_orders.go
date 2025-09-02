package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Pedidos struct {
	ID        int
	ClienteID int
	total     float64
}
type ItemDetalhe struct {
	Nome          string  `json:"nome"`
	Quantidade    int     `json:"quantidade"`
	PrecoUnitario float64 `json:"preco_unitario"`
}

type PedidoDetalhe struct {
	PedidoID  int           `json:"pedido_id"`
	ClienteID int           `json:"cliente_id"`
	Itens     []ItemDetalhe `json:"itens"`
	Total     float64       `json:"total"`
}

func list_pedidos(db *pgx.Conn) ([]Pedidos, error) {
	rows, err := db.Query(context.Background(),
		`SELECT 
			p.id AS pedido_id,
			p.cliente_id,
			SUM(i.preco_unitario * i.quantidade) AS total
		FROM pedidos p
		JOIN itens_pedido i ON i.pedido_id = p.id
		GROUP BY p.id, p.cliente_id
		ORDER BY p.id DESC
		LIMIT 10;`)
	if err != nil {
		return nil, fmt.Errorf("Error:", err)
	}
	defer rows.Close()

	var pedidos []Pedidos
	for rows.Next() {
		var p Pedidos
		if err := rows.Scan(&p.ID, &p.ClienteID, &p.total); err != nil {
			return nil, err
		}
		pedidos = append(pedidos, p)
	}
	return pedidos, nil
}

func get_pedido(db *pgx.Conn, pedidoID int) (*PedidoDetalhe, error) {
	// 1. Buscar cliente_id do pedido
	var clienteID int
	err := db.QueryRow(context.Background(),
		`SELECT cliente_id FROM pedidos WHERE id = $1`, pedidoID).Scan(&clienteID)
	if err != nil {
		return nil, fmt.Errorf("pedido não encontrado: %w", err)
	}

	// 2. Buscar itens do pedido
	rows, err := db.Query(context.Background(),
		`SELECT p.nome, i.quantidade, i.preco_unitario
		 FROM itens_pedido i
		 JOIN produtos p ON i.produto_id = p.id
		 WHERE i.pedido_id = $1`, pedidoID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar itens: %w", err)
	}
	defer rows.Close()

	var itens []ItemDetalhe
	var total float64
	for rows.Next() {
		var item ItemDetalhe
		if err := rows.Scan(&item.Nome, &item.Quantidade, &item.PrecoUnitario); err != nil {
			return nil, err
		}
		total += float64(item.Quantidade) * item.PrecoUnitario
		itens = append(itens, item)
	}

	pedido := &PedidoDetalhe{
		PedidoID:  pedidoID,
		ClienteID: clienteID,
		Itens:     itens,
		Total:     total,
	}

	return pedido, nil
}
