package main

import (
	"net/http"
	"simpleapi_go/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/api/clientes/iniciar_sesion", initSession)
	router.POST("/api/compras", createPurchase)
	router.GET("api/productos", getProducts)
	router.PUT("api/productos/:id", updateProduct)
	router.POST("api/productos/", createProduct)
	router.DELETE("api/productos/:id", deleteProduct)
	router.GET("api/estadisticas", getStatistics)
	router.Run(":5000")
}

// Funcionando la obtención de productos
func getProducts(c *gin.Context) {
	products := models.GetProducts()
	if products == nil || len(products) == 0 {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, products)
	}
}

// Funcionando la creación de productos
func createProduct(c *gin.Context) {
	var product models.Product
	c.BindJSON(&product)
	id := models.CreateProduct(product)
	c.JSON(http.StatusOK, gin.H{"id_producto": id})
}

// Funcionando la eliminación de productos
func deleteProduct(c *gin.Context) {
	id := c.Params.ByName("id")
	idInt, _ := strconv.Atoi(id)
	models.DeleteProduct(idInt)
	c.JSON(http.StatusOK, gin.H{"id_producto": idInt})
}

// Funcionando la actualización de productos
func updateProduct(c *gin.Context) {
	id := c.Params.ByName("id")
	idInt, _ := strconv.Atoi(id)
	var product models.Product
	c.BindJSON(&product)
	models.UpdateProduct(idInt, product)
	c.JSON(http.StatusOK, gin.H{"id_producto": idInt})
}

// Funcionando la verificación de usuario en sistema
func initSession(c *gin.Context) {
	var client models.Cliente
	c.BindJSON(&client)
	valid := models.InitSession(client)
	c.JSON(http.StatusOK, gin.H{"acceso_valido": valid})
}

// Funcionando la creación de compra con el detalle respectivo, incluido el timestamp
func createPurchase(c *gin.Context) {
	var compra models.BodyCompra
	c.BindJSON(&compra)
	id_compra, id_despacho, suma_total := models.CreatePurchase(compra)
	if id_compra == -1 {
		c.AbortWithStatus(http.StatusNotFound)
	}
	c.JSON(http.StatusOK, gin.H{"id_compra": id_compra, "id_despacho": id_despacho, "monto_total": suma_total})
}

// Funcionando la obtención de estadísticas de los productos vendidos
func getStatistics(c *gin.Context) {
	statistics := models.GetStatistics()
	if statistics == nil || len(statistics) == 0 {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, gin.H{"producto_mas_vendido": statistics[0], "producto_menos_vendido": statistics[1], "producto_mayor_ganancia": statistics[2], "producto_menor_ganancia": statistics[3]})
	}

}
