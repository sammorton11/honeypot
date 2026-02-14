package models

type Attempt struct {
	Address string `json:"address"`
	Network string `json:"network"`
	Message string `json:"message"`
}
