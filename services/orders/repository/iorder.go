package repository

type IOrders interface {
	Get(id int) (Order, error)
	Add(Order) (Order, error)
	Update(Order) (Order, error)
	Delete(Order) error
}
