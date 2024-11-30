// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package models

const TableNameUser = "users"

// User mapped from table <users>
type User struct {
	UserID              int64       `gorm:"column:user_id;primaryKey" json:"user_id"`
	UserRoleID          int64       `gorm:"column:user_role_id;not null;default:1" json:"user_role_id"`
	Email               string      `gorm:"column:email;not null" json:"email"`
	Password            []uint8     `gorm:"column:password;not null" json:"password"`
	RefreshTokenVersion int64       `gorm:"column:refresh_token_version;not null" json:"refresh_token_version"`
	Balance             float32     `gorm:"column:balance;not null;default:0" json:"balance"`
	UserRole            UserRole    `json:"user_role"`
	PcOrders            []PcOrder   `json:"pc_orders"`
	DishOrders          []DishOrder `json:"dish_orders"`
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}
