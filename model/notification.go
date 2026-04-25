package model

type Notification struct {
	ID              int    `json:"id"`
	Chamado_id      string `json:"chamado_id" binding:"required"`
	Tipo            string `json:"tipo" binding:"required"`
	Cpf             string `json:"cpf,omitempty" binding:"required"`
	Cpf_encrypted   string `json:"-"`
	Cpf_bindex      string `json:"-"`
	Status_anterior string `json:"status_anterior" binding:"required"`
	Status_novo     string `json:"status_novo" binding:"required"`
	Titulo          string `json:"titulo" binding:"required"`
	Descricao       string `json:"descricao" binding:"required"`
	Timestamp       string `json:"timestamp" binding:"required"`
	Is_read         bool   `json:"is_read"`
}
