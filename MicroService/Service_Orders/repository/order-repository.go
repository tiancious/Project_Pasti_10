package repository

import (
	"strconv"

	"github.com/tiancious/Service_Orders/entity"
	"gorm.io/gorm"
)

type OrderRepository interface {
	InsertOrder(order entity.Order) entity.Order
	UpdateOrder(order entity.Order) entity.Order
	All(userID uint64) []entity.Order
	FindByID(OrderID uint64) entity.Order
	DeleteOrder(order entity.Order)
	AppoveOrder(order entity.Order) entity.Order
	RejectOrder(order entity.Order) entity.Order
	CancelOrder(order entity.Order) entity.Order
}

type orderConnection struct {
	connection *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderConnection{
		connection: db,
	}
}

func (db *orderConnection) InsertOrder(order entity.Order) entity.Order {
    var existingOrdersCount int64
    db.connection.Model(&entity.Order{}).Count(&existingOrdersCount)
    
    // Generate new order ID based on existing orders count
    newOrderID := existingOrdersCount + 1
    
    // Generate order code by combining latest order id with "ORD-" prefix
    order.Code = "ORD-" + strconv.FormatInt(newOrderID, 10)
    order.ID = uint64(newOrderID)

    db.connection.Save(&order)
    return order
}


func (db *orderConnection) UpdateOrder(order entity.Order) entity.Order {
	db.connection.Save(&order)
	return order
}

func (db *orderConnection) All(userID uint64) []entity.Order {
	var orders []entity.Order
	db.connection.Preload("OrderItems").Where("user_id = ?", userID).Find(&orders)
	return orders
}

func (db *orderConnection) FindByID(orderID uint64) entity.Order {
	var order entity.Order
	// preload order items
	db.connection.Preload("OrderItems").Where("id = ?", orderID).Find(&order)
	return order
}

func (db *orderConnection) DeleteOrder(order entity.Order) {
	db.connection.Delete(&order)
}

func (db *orderConnection) AppoveOrder(order entity.Order) entity.Order {
	// update order status to "approved"
	order.Status = "approved"
	db.connection.Save(&order)
	return order
}

func (db *orderConnection) RejectOrder(order entity.Order) entity.Order {
	// update order status to "rejected"
	db.connection.Save(&order)
	return order
}

func (db *orderConnection) CancelOrder(order entity.Order) entity.Order {
	db.connection.Save(&order)
	return order
}
