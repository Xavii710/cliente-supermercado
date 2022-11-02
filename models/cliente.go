package models

type Cliente struct {
	Id_cliente int    `json:"id_cliente" autoincrement:"true" primary_key:"true"`
	Nombre     string `json:"nombre"`
	Contrasena string `json:"contrasena"`
}
