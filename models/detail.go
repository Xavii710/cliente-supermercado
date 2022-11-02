package models

type BodyCompra struct {
	Id_cliente int `json:"id_cliente"`
	Productos  []struct {
		Id_producto int `json:"id_producto"`
		Cantidad    int `json:"cantidad"`
	} `json:"productos"`
}
