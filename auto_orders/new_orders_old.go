package main

import (
	"context"
	"fmt"
	"log"

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

func new_order(db *pgx.Conn) error {
	var clienteID int
	err := db.QueryRow(context.Background(), `SELECT id FROM clientes ORDER BY random() LIMIT 1`).Scan(&clienteID)
	if err != nil {
		return fmt.Errorf("erro ao buscar cliente: %w", err)
	}

	rows, err := db.Query(context.Background(), `select id, valor from produtos order by random() limit 5`)
	if err != nil {
		fmt.Println("Erro ao buscar o produto %s", err)
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
	err = db.QueryRow(context.Background(), `insert into pedidos (cliente_id) values ($1) returning id`, clienteID).Scan(&pedidoID)
	if err != nil {
		log.Fatal(err)
	}
	var itensPedidos []ItensPedidos
	for _, p := range produtos {
		qnt := 1 // ou gerar aleatório, ex: rand.Intn(3)+1

		_, err := db.Exec(context.Background(),
			`INSERT INTO itens_pedido (pedido_id, produto_id, quantidade, preco_unitario)
         VALUES ($1, $2, $3, $4)`,
			pedidoID, p.ID, qnt, p.Preco,
		)
		if err != nil {
			return fmt.Errorf("erro ao inserir item do pedido: %w", err)
		}

		itensPedidos = append(itensPedidos, ItensPedidos{
			PedidoID:      pedidoID,
			ProdutoID:     p.ID,
			Quantidade:    qnt,
			PrecoUnitario: p.Preco,
		})
	}
	return nil

}

func main() {
	totalPedidos := 10000
	for i := 0; i < totalPedidos; i++ {

		conn := connectDB()
		if err := new_order(conn); err != nil {
			log.Fatal("erro ao criar pedido %v", err)
		}
		conn.Close(context.Background())
	}
}
