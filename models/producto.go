package models

type Product struct {
	Id_producto         int    `json:"id_producto" autoincrement:"true" primary_key:"true"`
	Nombre              string `json:"nombre"`
	Cantidad_disponible int    `json:"cantidad_disponible"`
	Precio_unitario     int    `json:"precio_unitario"`
}

type Producto struct {
	Id_producto int `json:"id_producto"`
	Cantidad    int `json:"cantidad"`
}
