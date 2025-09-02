package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type order struct {
	ID    int
	preco float64
}

type Produto struct {
	ID    int
	Preco float64
}
type ItensPedidos struct {
	PedidoID      int
	ProdutoID     int
	Quantidade    int
	PrecoUnitario float64
}
type PedidoResumo struct {
	PedidoID  int     `json:"pedido_id"`
	ClienteID int     `json:"cliente_id"`
	Total     float64 `json:"total"`
}

func insertItensPedidos(db *pgx.Conn, pedidoID int, produtos []Produto) error {
	// Cria as linhas que serão inseridas
	rows := make([][]interface{}, len(produtos))
	for i, p := range produtos {
		qnt := 1 // ou rand.Intn(3)+1
		rows[i] = []interface{}{pedidoID, p.ID, qnt, p.Preco}
	}

	// CopyFrom exige uma estrutura Columnar e TableName
	_, err := db.CopyFrom(
		context.Background(),
		pgx.Identifier{"itens_pedido"}, // tabela
		[]string{"pedido_id", "produto_id", "quantidade", "preco_unitario"}, // colunas
		pgx.CopyFromRows(rows), // dados
	)
	if err != nil {
		return fmt.Errorf("erro ao inserir itens via copy: %w", err)
	}
	return nil
}
func new_order(db *pgx.Conn) error {
	// inicia a transação
	tx, err := db.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	// garante rollback caso dê erro
	defer tx.Rollback(context.Background())

	var clienteID int
	err = tx.QueryRow(context.Background(), `SELECT id FROM clientes ORDER BY random() LIMIT 1`).Scan(&clienteID)
	if err != nil {
		return fmt.Errorf("erro ao buscar cliente: %w", err)
	}

	rows, err := tx.Query(context.Background(), `select id, valor from produtos order by random() limit 5`)
	if err != nil {
		return fmt.Errorf("erro ao buscar produtos: %w", err)
	}
	defer rows.Close()

	var produtos []Produto
	for rows.Next() {
		var p Produto
		if err := rows.Scan(&p.ID, &p.Preco); err != nil {
			return err
		}
		produtos = append(produtos, p)
	}

	var pedidoID int
	err = tx.QueryRow(context.Background(), `insert into pedidos (cliente_id) values ($1) returning id`, clienteID).Scan(&pedidoID)
	if err != nil {
		return fmt.Errorf("erro ao criar pedido: %w", err)
	}

	// insert itens via COPY dentro da transação
	rowsCopy := make([][]interface{}, len(produtos))
	for i, p := range produtos {
		qnt := 1
		rowsCopy[i] = []interface{}{pedidoID, p.ID, qnt, p.Preco}
	}
	_, err = tx.CopyFrom(
		context.Background(),
		pgx.Identifier{"itens_pedido"},
		[]string{"pedido_id", "produto_id", "quantidade", "preco_unitario"},
		pgx.CopyFromRows(rowsCopy),
	)
	if err != nil {
		return fmt.Errorf("erro ao inserir itens via copy: %w", err)
	}

	// commit no final
	if err := tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("erro ao dar commit: %w", err)
	}

	return nil
}
func make_pedidos(conn *pgx.Conn, totalPedidos int) ([]PedidoResumo, error) {
	resumo := []PedidoResumo{}

	for i := 0; i < totalPedidos; i++ {
		// inicia transação
		tx, err := conn.Begin(context.Background())
		if err != nil {
			return nil, fmt.Errorf("erro ao iniciar transação: %w", err)
		}
		defer tx.Rollback(context.Background())

		var clienteID int
		err = tx.QueryRow(context.Background(), `SELECT id FROM clientes ORDER BY random() LIMIT 1`).Scan(&clienteID)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar cliente: %w", err)
		}

		rows, err := tx.Query(context.Background(), `SELECT id, valor FROM produtos ORDER BY random() LIMIT 5`)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar produtos: %w", err)
		}
		defer rows.Close()

		var produtos []Produto
		var total float64
		for rows.Next() {
			var p Produto
			if err := rows.Scan(&p.ID, &p.Preco); err != nil {
				return nil, err
			}
			produtos = append(produtos, p)
			total += p.Preco
		}

		var pedidoID int
		err = tx.QueryRow(context.Background(), `INSERT INTO pedidos (cliente_id) VALUES ($1) RETURNING id`, clienteID).Scan(&pedidoID)
		if err != nil {
			return nil, fmt.Errorf("erro ao criar pedido: %w", err)
		}

		// insert itens via COPY
		rowsCopy := make([][]interface{}, len(produtos))
		for i, p := range produtos {
			qnt := 1
			rowsCopy[i] = []interface{}{pedidoID, p.ID, qnt, p.Preco}
		}
		_, err = tx.CopyFrom(
			context.Background(),
			pgx.Identifier{"itens_pedido"},
			[]string{"pedido_id", "produto_id", "quantidade", "preco_unitario"},
			pgx.CopyFromRows(rowsCopy),
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao inserir itens via copy: %w", err)
		}

		if err := tx.Commit(context.Background()); err != nil {
			return nil, fmt.Errorf("erro ao dar commit: %w", err)
		}

		resumo = append(resumo, PedidoResumo{
			PedidoID:  pedidoID,
			ClienteID: clienteID,
			Total:     total,
		})
	}

	return resumo, nil
}
