package models

import "time"

type Detalle struct {
	Id_compra   int       `json:"id_compra" primary_key:"true"`
	Id_producto int       `json:"id_producto" primary_key:"true"`
	Cantidad    int       `json:"cantidad"`
	Fecha       time.Time `json:"fecha"`
}
