package model

type Receipt struct {
	OrderNumber int // Incrementing
	Coffee      *Coffee
	Err         error
}
