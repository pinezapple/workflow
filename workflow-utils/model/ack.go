package model

type Ack struct{
	Ack     string `json:"ack" gorm:"column:ack; primaryKey; autoIncrement:false"`
	Payload  []byte `json:"payload" gorm:"column:payload"`
}
