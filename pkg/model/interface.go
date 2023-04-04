package model

type IOrderMiddleware interface {
	OrderTaken(*Order)
	OrderCompleted(*Order)
}
