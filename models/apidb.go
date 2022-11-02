package models

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Credenciales para conectar con MySQL
const dbuser = "admin"
const dbpass = "12345678"
const dbname = "tarea_1_sd"

func InitSession(cliente Cliente) bool {
	// Open conection
	db, err := sql.Open("mysql", dbuser+":"+dbpass+"@tcp(127.0.0.1:3306)/"+dbname)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	// Verify existence of user
	results, err := db.Query("SELECT * FROM cliente WHERE id_cliente = ? AND contrasena = ?", cliente.Id_cliente, cliente.Contrasena)
	if err != nil {
		panic(err.Error())
	}
	for results.Next() {
		// Exist
		fmt.Println("true")
		return true
	}
	// Not exist
	return false
}

func CreatePurchase(Compra BodyCompra) (int, int, int) {
	// Open conection
	db, err := sql.Open("mysql", dbuser+":"+dbpass+"@tcp(127.0.0.1:3306)/"+dbname)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	// Insert compra
	insert, err := db.Query("INSERT INTO compra VALUES (NULL,?)", Compra.Id_cliente)
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
	// Obtain id de compra
	results, err := db.Query("SELECT id_compra FROM compra WHERE id_cliente = ? ORDER BY id_compra DESC LIMIT 1", Compra.Id_cliente)
	if err != nil {
		panic(err.Error())
	}
	defer results.Close()
	for results.Next() {
		var id int
		err = results.Scan(&id)
		if err != nil {
			panic(err.Error())
		}
		rand.Seed(time.Now().UnixNano())
		var despacho Despacho
		despacho.Id_despacho = rand.Intn(9*1000000-1*1000000+1) + 1*1000000
		despacho.Estado = "RECIBIDO"
		despacho.Id_compra = id
		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		if err != nil {
			panic(err.Error())
		}
		defer conn.Close()
		ch, err := conn.Channel()
		if err != nil {
			panic(err.Error())
		}
		defer ch.Close()

		err = ch.ExchangeDeclare(
			"logs",   // name
			"direct", // type
			true,     // durable
			false,    // auto-deleted
			false,    // internal
			false,    // no-wait
			nil,      // arguments
		)
		if err != nil {
			panic(err.Error())
		}

		q, err := ch.QueueDeclare(
			"despacho", // name
			true,       // durable
			false,      // delete when unused
			false,      // exclusive
			false,      // no-wait
			nil,        // arguments
		)
		if err != nil {
			panic(err.Error())
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		body, err := despacho.MarshalBinary()
		if err != nil {
			panic(err.Error())
		}
		err = ch.PublishWithContext(ctx,
			"logs", // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        body,
			})
		if err != nil {
			panic(err.Error())
		}

		// Insert detalle de compra
		var suma_total int
		for _, producto := range Compra.Productos {
			var product Product
			//obtain cantidad disponible
			results, err := db.Query("SELECT cantidad_disponible FROM producto WHERE id_producto = ?", producto.Id_producto)
			if err != nil {
				panic(err.Error())
			}
			defer results.Close()
			for results.Next() {
				err = results.Scan(&product.Cantidad_disponible)
				if err != nil {
					panic(err.Error())
				}
			}

			if producto.Cantidad <= product.Cantidad_disponible {
				insert, err := db.Query("INSERT INTO detalle VALUES (?,?,?,?)", id, producto.Id_producto, producto.Cantidad, time.Now())
				if err != nil {
					panic(err.Error())
				}
				defer insert.Close()
				// Update cantidad disponible
				update, err := db.Query("UPDATE producto SET cantidad_disponible = cantidad_disponible - ? WHERE id_producto = ?", producto.Cantidad, producto.Id_producto)
				if err != nil {
					panic(err.Error())
				}
				defer update.Close()
				// Obtain precio unitario
				results, err := db.Query("SELECT precio_unitario FROM producto WHERE id_producto = ?", producto.Id_producto)
				if err != nil {
					panic(err.Error())
				}
				defer results.Close()
				for results.Next() {
					err = results.Scan(&product.Precio_unitario)
					if err != nil {
						panic(err.Error())
					}
				}
				suma_total += product.Precio_unitario * producto.Cantidad
			}
		}
		return id, despacho.Id_despacho, suma_total
	}
	return -1, -1, -1
}

func GetProducts() []Product {
	// Open conection
	db, err := sql.Open("mysql", dbuser+":"+dbpass+"@tcp(127.0.0.1:3306)/"+dbname)
	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}
	defer db.Close()
	// Get products from table producto
	results, err := db.Query("SELECT * FROM producto")
	if err != nil {
		fmt.Println("Err", err.Error())
		return nil
	}
	products := []Product{}
	for results.Next() {
		var product Product
		err = results.Scan(&product.Id_producto, &product.Nombre, &product.Cantidad_disponible, &product.Precio_unitario)
		if err != nil {
			panic(err.Error())
		}
		products = append(products, product)
	}
	return products
}

func CreateProduct(product Product) int {
	// Open conection
	db, err := sql.Open("mysql", dbuser+":"+dbpass+"@tcp(127.0.0.1:3306)/"+dbname)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	// Insert product in db
	insert, err := db.Query("INSERT INTO producto VALUES (?,?,?,?)", product.Id_producto, product.Nombre, product.Cantidad_disponible, product.Precio_unitario)
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()

	// Obtain id of the product with same name, cantidad_disponible and precio_unitario
	results, err := db.Query("SELECT id_producto FROM producto WHERE nombre = ? AND cantidad_disponible = ? AND precio_unitario = ?", product.Nombre, product.Cantidad_disponible, product.Precio_unitario)
	if err != nil {
		panic(err.Error())
	}
	defer results.Close()
	for results.Next() {
		var id int
		err = results.Scan(&id)
		if err != nil {
			panic(err.Error())
		}
		return id
	}
	return 0
}

func UpdateProduct(id int, product Product) int {
	// Open conection
	db, err := sql.Open("mysql", dbuser+":"+dbpass+"@tcp(127.0.0.1:3306)/"+dbname)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	// Update product with id
	update, err := db.Query("UPDATE producto SET nombre = ?, cantidad_disponible = ?, precio_unitario = ? WHERE id_producto = ?", product.Nombre, product.Cantidad_disponible, product.Precio_unitario, id)
	if err != nil {
		panic(err.Error())
	}
	defer update.Close()
	return product.Id_producto
}

func DeleteProduct(id int) int {
	// Open connection
	db, err := sql.Open("mysql", dbuser+":"+dbpass+"@tcp(127.0.0.1:3306)/"+dbname)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	// Delete product with id
	delete, err := db.Query("DELETE FROM producto WHERE id_producto = ?", id)
	if err != nil {
		panic(err.Error())
	}
	defer delete.Close()
	return id
}

func GetStatistics() []int {
	// Create array resultados
	var resultados []int
	// Open connection
	db, err := sql.Open("mysql", dbuser+":"+dbpass+"@tcp(127.0.0.1:3306)/"+dbname)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	// Obtain most sold product from detalle
	results, err := db.Query("SELECT id_producto, SUM(cantidad) FROM detalle GROUP BY id_producto ORDER BY SUM(cantidad) DESC LIMIT 1")
	if err != nil {
		panic(err.Error())
	}
	defer results.Close()
	for results.Next() {
		var id int
		var cantidad int
		err = results.Scan(&id, &cantidad)
		if err != nil {
			panic(err.Error())
		}
		resultados = append(resultados, id)
	}
	// Obtain least sold product from detalle
	results, err = db.Query("SELECT id_producto, SUM(cantidad) FROM detalle GROUP BY id_producto ORDER BY SUM(cantidad) ASC LIMIT 1")
	if err != nil {
		panic(err.Error())
	}
	defer results.Close()
	for results.Next() {
		var id int
		var cantidad int
		err = results.Scan(&id, &cantidad)
		if err != nil {
			panic(err.Error())
		}
		resultados = append(resultados, id)
	}
	// Obtain highest profit from tables detalle and product
	results, err = db.Query("SELECT detalle.id_producto, SUM(detalle.cantidad) * producto.precio_unitario FROM detalle INNER JOIN producto ON detalle.id_producto = producto.id_producto GROUP BY detalle.id_producto ORDER BY SUM(detalle.cantidad) * producto.precio_unitario DESC LIMIT 1")
	if err != nil {
		fmt.Println("Pasó por acá")
		panic(err.Error())
	}
	defer results.Close()
	for results.Next() {
		var id int
		var ganancia int
		err = results.Scan(&id, &ganancia)
		if err != nil {
			panic(err.Error())
		}
		resultados = append(resultados, id)
	}
	// Obtain lowest profit from tables detalle and product
	results, err = db.Query("SELECT detalle.id_producto, SUM(detalle.cantidad) * producto.precio_unitario FROM detalle INNER JOIN producto ON detalle.id_producto = producto.id_producto GROUP BY detalle.id_producto ORDER BY SUM(detalle.cantidad) * producto.precio_unitario ASC LIMIT 1")
	if err != nil {
		panic(err.Error())
	}
	defer results.Close()
	for results.Next() {
		var id int
		var ganancia int
		err = results.Scan(&id, &ganancia)
		if err != nil {
			panic(err.Error())
		}
		resultados = append(resultados, id)
	}
	return resultados
}
