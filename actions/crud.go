package actions

import (
	"net/http"
	f "teste/models"

	"github.com/gobuffalo/buffalo"
)

type UserResponse struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	DataHoraCriacao string `json:"data_hora_criacao"`
	Password        string `json:"password"`
}

func SelectHandler(c buffalo.Context) error {
	teste := f.FindAll()
	response := []UserResponse{}

	for _, v := range teste {
		user := UserResponse{int(v.V_id.Int64), v.V_name.String, v.V_data_hora_criacao.String, v.V_password.String}
		response = append(response, user)
	}
	return c.Render(http.StatusOK, r.JSON(response))
}

func SelectIdHandler(c buffalo.Context) error {
	teste := f.FindId(c.Param("id"))
	response := UserResponse{int(teste[0].V_id.Int64), teste[0].V_name.String, teste[0].V_data_hora_criacao.String, teste[0].V_password.String}
	return c.Render(http.StatusOK, r.JSON(response))
}

func DeleteIdHandler(c buffalo.Context) error {
	teste := f.DeleteById(c.Param("id"))
	return c.Render(http.StatusOK, r.JSON(teste))
}

type UpdateRequest struct {
	Id   string `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

func UpdateIdHandler(c buffalo.Context) error {
	var req UpdateRequest
	c.Bind(&req)

	teste := f.UpdateById(req.Id, req.Name)
	response := UserResponse{int(teste[0].V_id.Int64), teste[0].V_name.String, teste[0].V_data_hora_criacao.String, teste[0].V_password.String}

	return c.Render(http.StatusOK, r.JSON(response))
}

type InsertRequest struct {
	Nome  string `json:"name" validate:"required,min=8,max=32"`
	Senha string `json:"password" validate:"required"`
}

func InsertHandler(c buffalo.Context) error {
	var req InsertRequest
	c.Bind(&req)

	teste := f.Insert(req.Nome, req.Senha)[0]
	response := UserResponse{int(teste.V_id.Int64), teste.V_name.String, teste.V_data_hora_criacao.String, teste.V_password.String}

	return c.Render(http.StatusOK, r.JSON(response))
}
