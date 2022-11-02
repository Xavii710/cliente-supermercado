package models

import (
	"encoding/json"
	"fmt"
)

type Despacho struct {
	Id_despacho int    `json:"id_despacho"`
	Estado      string `json:"estado"`
	Id_compra   int    `json:"id_compra"`
}

// MarshalBinary is a custom marshaler for Despacho
func (d *Despacho) MarshalBinary() ([]byte, error) {
	body, error := json.Marshal(d)
	if error != nil {
		fmt.Println(error)
	}
	return body, nil
}
