package models

import (
	"database/sql"
	"encoding/json"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
)

// UsersTbl is used by pop to map your users_tbls database table to your go code.
type UsersTbl struct {
	V_id                sql.NullInt64  `json:"id" db:"id"`
	V_name              sql.NullString `json:"name" db:"name"`
	V_data_hora_criacao sql.NullString `json:"data_hora_criacao" db:"data_hora_criacao"`
	V_password          sql.NullString `json:"password" db:"password"`
}

func FindAll() []UsersTbl {
	db_consulta, _ := pop.Connect("development")
	ListaRetorno := []UsersTbl{}

	db_consulta.RawQuery("SELECT * FROM users_tbl").All(&ListaRetorno)
	return ListaRetorno
}

func FindId(id string) []UsersTbl {
	db_consulta, _ := pop.Connect("development")
	retorno := []UsersTbl{}

	db_consulta.RawQuery("SELECT * FROM users_tbl WHERE id = ?", id).All(&retorno)
	return retorno
}

func DeleteById(id string) interface{} {
	db_consulta, _ := pop.Connect("development")
	retorno := []UsersTbl{}

	db_consulta.RawQuery("DELETE FROM users_tbl WHERE id = ?", id).Exec()
	db_consulta.RawQuery("SELECT * FROM users_tbl").All(&retorno)
	return "Removido com sucesso!"
}

func UpdateById(id string, name string) []UsersTbl {
	db_consulta, _ := pop.Connect("development")
	retorno := []UsersTbl{}

	db_consulta.RawQuery("UPDATE users_tbl SET name = ? WHERE id = ?", name, id).Exec()
	db_consulta.RawQuery("SELECT * FROM users_tbl WHERE id = ?", id).All(&retorno)
	return retorno
}

func Insert(name string, password string) []UsersTbl {
	db_consulta, _ := pop.Connect("development")
	retorno := []UsersTbl{}

	db_consulta.RawQuery("INSERT INTO users_tbl SET name=?, password=?", name, password).Exec()
	db_consulta.RawQuery("SELECT * FROM users_tbl ORDER BY id DESC LIMIT 1").All(&retorno)
	return retorno
}

func AuthUserByUsername(username string) []UsersTbl {
	db_consulta, _ := pop.Connect("development")
	User := []UsersTbl{}

	db_consulta.RawQuery("SELECT * FROM users_tbl WHERE name = ?", username).All(&User)
	return User
}

func AuthUserById(id string) []UsersTbl {
	db_consulta, _ := pop.Connect("development")
	User := []UsersTbl{}

	db_consulta.RawQuery("SELECT * FROM users_tbl WHERE id = ?", id).All(&User)
	return User
}

// String is not required by pop and may be deleted
func (u UsersTbl) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// UsersTbls is not required by pop and may be deleted
type UsersTbls []UsersTbl

// String is not required by pop and may be deleted
func (u UsersTbls) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *UsersTbl) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *UsersTbl) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *UsersTbl) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
