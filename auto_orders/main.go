package main

import (
	"context"
	"net/http"
	"strconv"

	_ "autoorders/rdls.dev/docs"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// PingHandler verifica se o serviço está ativo
// @Summary Verifica se o serviço está ativo
// @Description Retorna "pong" para testar a API
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /ping [get]
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

// PongHandler retorna "ping"
// @Summary Teste alternativo
// @Description Retorna "ping"
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /pong [get]
func PongHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ping"})
}

// MakePedidosGetHandler cria 10 pedidos aleatórios via GET
// @Summary Criar pedidos aleatórios
// @Description Gera 10 pedidos aleatórios
// @Tags Pedidos
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /make_pedidos [get]
func MakePedidosGetHandler(c *gin.Context) {
	conn := connectDB()
	defer conn.Close(context.Background())

	totalPedidos := 10
	resultado, err := make_pedidos(conn, totalPedidos)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": resultado})
}

// MakePedidosPostHandler cria N pedidos aleatórios via POST
// @Summary Criar pedidos via POST
// @Description Gera N pedidos aleatórios informando no body
// @Tags Pedidos
// @Accept json
// @Produce json
// @Param total_pedidos body int true "Número de pedidos a criar"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /make_pedidos [post]
func MakePedidosPostHandler(c *gin.Context) {
	var body struct {
		TotalPedidos int `json:"total_pedidos"`
	}

	if err := c.ShouldBindJSON(&body); err != nil || body.TotalPedidos <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "informe um total_pedidos válido"})
		return
	}

	conn := connectDB()
	defer conn.Close(context.Background())

	resumo, err := make_pedidos(conn, body.TotalPedidos)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pedidos": resumo})
}

// GetPedidoHandler retorna os itens de um pedido pelo ID
// @Summary Detalhes de um pedido
// @Description Retorna os itens de um pedido pelo ID
// @Tags Pedidos
// @Produce json
// @Param id path int true "ID do pedido"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /pedido/{id} [get]
func GetPedidoHandler(c *gin.Context) {
	idStr := c.Param("id")
	pedidoID, err := strconv.Atoi(idStr)
	if err != nil || pedidoID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID inválido"})
		return
	}

	conn := connectDB()
	defer conn.Close(context.Background())

	pedido, err := get_pedido(conn, pedidoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pedido)
}

// ListPedidosHandler lista os últimos 10 pedidos
// @Summary Últimos pedidos
// @Description Retorna os últimos 10 pedidos com total e cliente
// @Tags Pedidos
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /ultimos_pedidos [get]
func ListPedidosHandler(c *gin.Context) {
	conn := connectDB()
	defer conn.Close(context.Background())

	pedidos, err := list_pedidos(conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pedidos": pedidos})
}

func main() {
	r := gin.Default()

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Rotas
	r.GET("/ping", PingHandler)
	r.GET("/pong", PongHandler)
	r.GET("/make_pedidos", MakePedidosGetHandler)
	r.POST("/make_pedidos", MakePedidosPostHandler)
	r.GET("/pedido/:id", GetPedidoHandler)
	r.GET("/ultimos_pedidos", ListPedidosHandler)

	r.Run()
}
