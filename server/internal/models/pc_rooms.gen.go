// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package models

const TableNamePcRoom = "pc_rooms"

// PcRoom mapped from table <pc_rooms>
type PcRoom struct {
	PcRoomID    int64  `gorm:"column:pc_room_id;primaryKey" json:"pc_room_id"`
	Name        string `gorm:"column:name;not null" json:"name"`
	Rows        int    `gorm:"column:rows;not null" json:"rows"`
	Places      int    `gorm:"column:places;not null" json:"places"`
	Description string `gorm:"column:description" json:"description"`
	Pcs         []Pc   `json:"pcs"`
}

// TableName PcRoom's table name
func (*PcRoom) TableName() string {
	return TableNamePcRoom
}
