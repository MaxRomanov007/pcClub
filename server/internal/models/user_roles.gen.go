// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package models

const TableNameUserRole = "user_roles"

// UserRole mapped from table <user_roles>
type UserRole struct {
	UserRoleID int64  `gorm:"column:user_role_id;primaryKey" json:"user_role_id"`
	Name       string `gorm:"column:name;not null" json:"name"`
	Users      []User `json:"users"`
}

// TableName UserRole's table name
func (*UserRole) TableName() string {
	return TableNameUserRole
}
