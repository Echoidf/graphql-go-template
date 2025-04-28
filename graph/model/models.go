package model

type Order struct {
	Id           string `json:"id"`
	InstrumentId string `json:"instrumentId"`
	OrderId      string `json:"orderId"`
}
