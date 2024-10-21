package models

import "database/sql"

// PcType represents the pc_types table.
type PcType struct {
	PcTypeID    int64  `db:"pc_type_id"`
	ProcessorID int64  `db:"processor_id"`
	VideoCardID int64  `db:"video_card_id"`
	MonitorID   int64  `db:"monitor_id"`
	RamID       int64  `db:"ram_id"`
	Name        string `db:"name"`

	Processor *Processor
	VideoCard *VideoCard
	Monitor   *Monitor
	Ram       *Ram
}

type PcTypeData struct {
	TypeID      int64          `json:"id"`
	TypeName    string         `json:"name"`
	Description string         `json:"description"`
	Processor   *ProcessorData `json:"processor"`
	VideoCard   *VideoCardData `json:"video_card"`
	Monitor     *MonitorData   `json:"monitor"`
	Ram         *RamData       `json:"ram"`
}

// PcTypeImage represents the pc_type_image table.
type PcTypeImage struct {
	PcTypeImageID int64  `db:"pc_type_image_id"`
	PcTypeID      int64  `db:"pc_type_id"`
	IsMain        bool   `db:"is_main"`
	Path          string `db:"path"`

	PcType *PcType
}

// PcRoom represents the pc_rooms table.
type PcRoom struct {
	PcRoomID    int64  `db:"pc_room_id"`
	Name        string `db:"name"`
	Rows        int    `db:"rows"`
	Places      int    `db:"places"`
	Description string `db:"description"`
}

// PcStatus represents the pc_statuses table.
type PcStatus struct {
	PcStatusID int64  `db:"pc_status_id"`
	Name       string `db:"name"`
}

// Pc represents the pc table.
type Pc struct {
	PcID        int64            `db:"pc_id"`
	PcRoomID    int64            `db:"pc_room_id"`
	PcTypeID    int64            `db:"pc_type_id"`
	PcStatusID  int64            `db:"pc_status_id"`
	Row         int              `db:"row"`
	Place       int              `db:"place"`
	Description sql.Null[string] `db:"description"`
	PcRoom      *PcRoom
	PcType      *PcType
	PcStatus    *PcStatus
}

// PcData represents the pc data
type PcData struct {
	PcID        int64  `db:"pc_id" json:"pc_id"`
	Row         int    `db:"row" json:"row"`
	Place       int    `db:"place" json:"place"`
	Description string `db:"description" json:"description"`
	PcRoomID    int64  `db:"pc_room_id" json:"pc_room_id"`
}

// PcOrderStatus represents the pc_order_statuses table.
type PcOrderStatus struct {
	PcOrderStatusID int64  `db:"pc_order_status_id"`
	Name            string `db:"name"`
}

// PcOrder represents the pc_orders table.
type PcOrder struct {
	PcOrderID       int64   `db:"pc_order_id"`
	UserID          int64   `db:"user_id"`
	PcID            int64   `db:"pc_id"`
	PcOrderStatusID int64   `db:"pc_order_status_id"`
	Code            string  `db:"code"`
	Cost            float64 `db:"cost"`
	StartTime       string  `db:"start_time"`
	EndTime         string  `db:"end_time"`
	ActualEndTime   string  `db:"actual_end_time"`
	OrderDate       string  `db:"order_date"`
	User            *User
	Pc              *Pc
	PcOrderStatus   *PcOrderStatus
}
