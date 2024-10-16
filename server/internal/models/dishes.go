package models

import "time"

// DishStatus represents the dish_statuses table.
type DishStatus struct {
	DishStatusID int64  `db:"dish_status_id"`
	Name         string `db:"name"`
}

// Dish represents the dishes table.
type Dish struct {
	DishID       int64   `db:"dish_id"`
	DishStatusID int64   `db:"dish_status_id"`
	Name         string  `db:"name"`
	Calories     int16   `db:"calories"`
	Cost         float64 `db:"cost"`
	Description  string  `db:"description"`
	DishStatus   *DishStatus
}

// DishImage represents the dish_images table.
type DishImage struct {
	DishImageID int64  `db:"dish_image_id"`
	DishID      int64  `db:"dish_id"`
	IsMain      bool   `db:"is_main"`
	Path        string `db:"path"`
	Dish        *Dish
}

// DishOrderStatus represents the dish_order_statuses table.
type DishOrderStatus struct {
	DishOrderStatusID int64  `db:"dish_order_status_id"`
	Name              string `db:"name"`
}

// DishOrder represents the dish_orders table.
type DishOrder struct {
	DishOrderID       int64     `db:"dish_order_id"`
	DishID            int64     `db:"dish_id"`
	UserID            int64     `db:"user_id"`
	DishOrderStatusID int64     `db:"dish_order_status_id"`
	Count             int16     `db:"count"`
	Cost              float64   `db:"cost"`
	OrderDate         time.Time `db:"order_date"`
	Dish              *Dish
	User              *User
	DishOrderStatus   *DishOrderStatus
}
