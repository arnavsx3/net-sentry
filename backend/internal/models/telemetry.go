package models

type TargetPayload struct {
	Host string `json:"host" binding:"required"`
}
