package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"simpleapi_go/models"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Bienvenido")
	// Label main to break the loop
main:
	for {
		fmt.Println("\nOpciones: \n1. Iniciar sesión como cliente \n2. Iniciar sesión como administrador \n3. Salir")
		var option int
		var compra_realizada bool
		compra_realizada = false
		fmt.Scanln(&option)
		switch option {
		case 1:
			fmt.Print("Ingrese id: ")
			var id int
			fmt.Scanln(&id)
			fmt.Print("Ingrese contraseña: ")
			var password string
			fmt.Scanln(&password)
			// jsonbody format example {"id_cliente": 1, "contrasena": "1234"}
			jsonBody := fmt.Sprintf(`{"id_cliente": %d, "contrasena": "%s"}`, id, password)
			body := []byte(jsonBody)
			fmt.Println(bytes.NewBuffer(body))
			//Post a la RESTAPI SERVER para iniciar sesión
			resp, err := http.NewRequest("POST", "http://localhost:5000/api/clientes/iniciar_sesion", bytes.NewBuffer(body))
			if err != nil {
				fmt.Println(err)
			}
			// Set the header
			resp.Header.Add("content-type", "application/x-www-form-urlencoded")
			client := &http.Client{}
			res, err := client.Do(resp)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()
			// Read the body
			new_body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				panic(err)
			}
			// Check if the login was successful with "true" string contained in the response body
			if strings.Contains(string(new_body), "true") {
				fmt.Println("Inicio de sesión exitoso")
				// Label login to break the loop

			login:
				for {
					if compra_realizada == false {
						fmt.Println("\nOpciones: \n1. Ver lista de productos \n2. Hacer Compra \n3. Salir")
						fmt.Print("Ingrese una opción: ")
						var opcion int
						fmt.Scanln(&opcion)
						switch opcion {
						case 1:
							// Get a la RESTAPI SERVER para obtener lista de productos
							resp, err := http.Get("http://localhost:5000/api/productos")
							if err != nil {
								fmt.Println(err)
							}
							defer resp.Body.Close()
							// Read the body
							body, err := ioutil.ReadAll(resp.Body)
							if err != nil {
								fmt.Println(err)
							}
							// Print products
							// string body to []byte
							var productos []models.Product
							json.Unmarshal(body, &productos)
							// Print products
							for _, producto := range productos {
								fmt.Printf("%d;%s;%d por unidad;%d disponibles\n", producto.Id_producto, producto.Nombre, producto.Precio_unitario, producto.Cantidad_disponible)
							}
						case 2:
							var cantidad_total int
							fmt.Print("Ingrese cantidad de productos a comprar: ")
							var cantidad int
							fmt.Scanln(&cantidad)
							// use BodyCompra struct
							var BodyCompra models.BodyCompra
							BodyCompra.Id_cliente = id
							for i := 0; i < cantidad; i++ {
								fmt.Printf("Ingrese producto %d par id-cantidad: ", i+1)
								// par id-cantidad input example: 1-2
								var par string
								fmt.Scanln(&par)
								//split par
								split := strings.Split(par, "-")
								//obtain id_producto
								id_producto, _ := strconv.Atoi(split[0])
								//obtain cantidad
								cantidad, _ := strconv.Atoi(split[1])
								cantidad_total += cantidad
								//append to BodyCompra
								BodyCompra.Productos = append(BodyCompra.Productos, models.Producto{Id_producto: id_producto, Cantidad: cantidad})
							}
							body, error := json.Marshal(BodyCompra)
							if error != nil {
								fmt.Println(error)
							}
							resp, err := http.NewRequest("POST", "http://localhost:5000/api/compras", bytes.NewBuffer(body))
							if err != nil {
								fmt.Println(err)
							}
							resp.Header.Add("content-type", "application/x-www-form-urlencoded")
							client := &http.Client{}
							res, err := client.Do(resp)
							if err != nil {
								panic(err)
							}
							defer res.Body.Close()
							new_body, err := ioutil.ReadAll(res.Body)
							if err != nil {
								panic(err)
							}
							//obtain json values from response body
							var compra models.CompraRealizada
							json.Unmarshal(new_body, &compra)
							//print compra
							fmt.Println("Gracias por su compra!")
							fmt.Printf("Cantidad de productos comprados: %d\n", cantidad_total)
							fmt.Printf("Monto total de la compra: %d\n", compra.Monto_total)
							fmt.Printf("El ID del despacho es %d\n", compra.Id_despacho)
							compra_realizada = true

						case 3:
							break login
						default:
							fmt.Println("Opción no válida")
						}
					} else {
						fmt.Println("\nOpciones: \n1. Ver lista de productos \n2. Hacer Compra \n3. Consultar despacho \n4. Salir")
						fmt.Print("Ingrese una opción: ")
						var opcion int
						fmt.Scanln(&opcion)
						switch opcion {
						case 1:
							// Get a la RESTAPI SERVER para obtener lista de productos
							resp, err := http.Get("http://localhost:5000/api/productos")
							if err != nil {
								fmt.Println(err)
							}
							defer resp.Body.Close()
							// Read the body
							body, err := ioutil.ReadAll(resp.Body)
							if err != nil {
								fmt.Println(err)
							}
							// Print products
							// string body to []byte
							var productos []models.Product
							json.Unmarshal(body, &productos)
							// Print products
							for _, producto := range productos {
								fmt.Printf("%d;%s;%d por unidad;%d disponibles\n", producto.Id_producto, producto.Nombre, producto.Precio_unitario, producto.Cantidad_disponible)
							}
						case 2:
							var cantidad_total int
							fmt.Print("Ingrese cantidad de productos a comprar: ")
							var cantidad int
							fmt.Scanln(&cantidad)
							// use BodyCompra struct
							var BodyCompra models.BodyCompra
							BodyCompra.Id_cliente = id
							for i := 0; i < cantidad; i++ {
								fmt.Printf("Ingrese producto %d par id-cantidad: ", i+1)
								// par id-cantidad input example: 1-2
								var par string
								fmt.Scanln(&par)
								//split par
								split := strings.Split(par, "-")
								//obtain id_producto
								id_producto, _ := strconv.Atoi(split[0])
								//obtain cantidad
								cantidad, _ := strconv.Atoi(split[1])
								cantidad_total += cantidad
								//append to BodyCompra
								BodyCompra.Productos = append(BodyCompra.Productos, models.Producto{Id_producto: id_producto, Cantidad: cantidad})
							}
							body, error := json.Marshal(BodyCompra)
							if error != nil {
								fmt.Println(error)
							}
							resp, err := http.NewRequest("POST", "http://localhost:5000/api/compras", bytes.NewBuffer(body))
							if err != nil {
								fmt.Println(err)
							}
							resp.Header.Add("content-type", "application/x-www-form-urlencoded")
							client := &http.Client{}
							res, err := client.Do(resp)
							if err != nil {
								panic(err)
							}
							defer res.Body.Close()
							new_body, err := ioutil.ReadAll(res.Body)
							if err != nil {
								panic(err)
							}
							//obtain json values from response body
							var compra models.CompraRealizada
							json.Unmarshal(new_body, &compra)
							//print compra
							fmt.Println("Gracias por su compra!")
							fmt.Printf("Cantidad de productos comprados: %d\n", cantidad_total)
							fmt.Printf("Monto total de la compra: %d\n", compra.Monto_total)
							fmt.Printf("El ID del despacho es %d\n", compra.Id_despacho)
						case 3:
							fmt.Print("Ingrese el ID del despacho: ")
							var id_despacho int
							fmt.Scanln(&id_despacho)
							resp, err := http.Get("http://localhost:5001/api/clientes/estado_despacho/" + strconv.Itoa(id_despacho))
							if err != nil {
								fmt.Println(err)
							}
							defer resp.Body.Close()
							// Read the body
							body, err := ioutil.ReadAll(resp.Body)
							if err != nil {
								fmt.Println(err)
							}
							// Print products
							// string body to []byte
							var estado_despacho models.Despacho
							json.Unmarshal(body, &estado_despacho)
							estado_despacho.Id_despacho = id_despacho
							// Print products
							fmt.Printf("Estado del despacho %d: %s\n", estado_despacho.Id_despacho, estado_despacho.Estado)

						case 4:
							break login
						}
					}
				}
			} else {
				fmt.Println("Inicio de sesión fallido")
			}

		case 2:
			//Enter admin password
			fmt.Println("Ingrese contraseña de administrador: ")
			var password string
			fmt.Scanln(&password)
			//if password is 1234
			if password == "1234" {
				fmt.Println("Inicio de sesión exitoso")
			adminlogin:
				for {
					fmt.Println("\nOpciones: \n1. Ver lista de productos \n2. Crear producto \n3. Eliminar producto \n4. Ver estadísticas \n5. Salir")
					fmt.Print("Ingrese una opción: ")
					var opcion int
					fmt.Scanln(&opcion)
					switch opcion {
					case 1:
						resp, err := http.Get("http://localhost:5000/api/productos")
						if err != nil {
							fmt.Println(err)
						}
						defer resp.Body.Close()
						body, err := ioutil.ReadAll(resp.Body)
						if err != nil {
							fmt.Println(err)
						}
						// Print products
						// string body to []byte
						var productos []models.Product
						json.Unmarshal(body, &productos)
						// Print products
						for _, producto := range productos {
							fmt.Printf("%d;%s;%d por unidad;%d disponibles\n", producto.Id_producto, producto.Nombre, producto.Precio_unitario, producto.Cantidad_disponible)
						}
					case 2:
						fmt.Println("Ingrese el nombre: ")
						var nombre string
						fmt.Scanln(&nombre)
						fmt.Println("Ingrese la disponibilidad: ")
						var disponibilidad int
						fmt.Scanln(&disponibilidad)
						fmt.Println("Ingrese el precio unitario: ")
						var precio_unitario int
						fmt.Scanln(&precio_unitario)
						// use Product struct
						var producto models.Product
						producto.Nombre = nombre
						producto.Cantidad_disponible = disponibilidad
						producto.Precio_unitario = precio_unitario
						body, error := json.Marshal(producto)
						if error != nil {
							fmt.Println(error)
						}
						// Post request
						resp, err := http.NewRequest("POST", "http://localhost:5000/api/productos", bytes.NewBuffer(body))
						if err != nil {
							fmt.Println(err)
						}
						resp.Header.Add("content-type", "application/x-www-form-urlencoded")
						client := &http.Client{}
						res, err := client.Do(resp)
						if err != nil {
							panic(err)
						} else {
							fmt.Println("Producto creado exitosamente")
						}
						defer res.Body.Close()
					case 3:
						fmt.Println("Ingrese el id del producto a eliminar: ")
						var id_producto int
						fmt.Scanln(&id_producto)
						// use Product struct
						var producto models.Product
						producto.Id_producto = id_producto
						body, error := json.Marshal(producto)
						if error != nil {
							fmt.Println(error)
						}
						// Delete request
						var url = fmt.Sprintf("http://localhost:5000/api/productos/%d", id_producto)
						resp, err := http.NewRequest("DELETE", url, bytes.NewBuffer(body))

						if err != nil {
							fmt.Println(err)
						}
						resp.Header.Add("content-type", "application/x-www-form-urlencoded")
						client := &http.Client{}
						res, err := client.Do(resp)
						if err != nil {
							panic(err)
						} else {
							fmt.Println("Producto eliminado exitosamente")
						}
						defer res.Body.Close()
					case 4:
						resp, err := http.Get("http://localhost:5000/api/estadisticas")
						if err != nil {
							fmt.Println(err)
						}
						defer resp.Body.Close()
						body, err := ioutil.ReadAll(resp.Body)
						if err != nil {
							fmt.Println(err)
						}
						fmt.Println(string(body))
					case 5:
						break adminlogin
					}
				}
			} else {
				fmt.Println("Inicio de sesión fallido")
			}
		case 3:
			fmt.Println("Hasta luego!")
			break main
		}
	}
}
