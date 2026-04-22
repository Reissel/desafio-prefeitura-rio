package model

type Notification struct {
	ID              int    `json:"id"`
	Chamado_id      string `json:"chamado_id"`
	Tipo            string `json:"tipo"`
	Cpf             string `json:"cpf"`
	Cpf_encrypted   string `json:"-"`
	Cpf_bindex      string `json:"-"`
	Status_anterior string `json:"status_anterior"`
	Status_novo     string `json:"status_novo"`
	Titulo          string `json:"titulo"`
	Descricao       string `json:"descricao"`
	Timestamp       string `json:"timestamp"`
	Is_read         bool   `json:"is_read"`
}
