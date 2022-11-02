package models

type Compra struct {
	Id_compra  int `json:"id_compra" autoincrement:"true" primary_key:"true"`
	Id_cliente int `json:"id_cliente"`
}

type CompraRealizada struct {
	Id_compra   int `json:"id_compra"`
	Id_despacho int `json:"id_despacho"`
	Monto_total int `json:"monto_total"`
}
