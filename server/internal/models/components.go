package models

// ProcessorProducer represents the processor_producers table.
type ProcessorProducer struct {
	ProcessorProducerID int64  `db:"processor_producer_id" json:"id"`
	Name                string `db:"name" json:"name"`
}

// VideoCardProducer represents the video_card_producers table.
type VideoCardProducer struct {
	VideoCardProducerID int64  `db:"video_card_producer_id" json:"id"`
	Name                string `db:"name" json:"name"`
}

// MonitorProducer represents the monitor_producers table.
type MonitorProducer struct {
	MonitorProducerID int64  `db:"monitor_producer_id" json:"id"`
	Name              string `db:"name" json:"name"`
}

// RamType represents the ram_types table.
type RamType struct {
	RamTypeID int64  `db:"ram_type_id" json:"id"`
	Name      string `db:"name" json:"name"`
}

// Processor represents the processors table.
type Processor struct {
	ProcessorID         int64  `db:"processor_id" json:"id"`
	ProcessorProducerID int64  `db:"processor_producer_id" json:"producer_id"`
	Model               string `db:"model" json:"model"`
}

// VideoCard represents the video_cards table.
type VideoCard struct {
	VideoCardID         int64  `db:"video_card_id" json:"id"`
	VideoCardProducerID int64  `db:"video_card_producer_id" json:"producer_id"`
	Model               string `db:"model" json:"model"`
}

// Monitor represents the monitors table.
type Monitor struct {
	MonitorID         int64  `db:"monitor_id" json:"id"`
	MonitorProducerID int64  `db:"monitor_producer_id" json:"producer_id"`
	Model             string `db:"model" json:"model"`
}

// Ram represents the ram table.
type Ram struct {
	RamID     int64 `db:"ram_id" json:"id"`
	RamTypeID int64 `db:"ram_type_id" json:"type_id"`
	Capacity  int   `db:"capacity" json:"capacity"`
}

type ProcessorData struct {
	Model    string `json:"model" validate:"required,max=255"`
	Producer string `json:"producer" validate:"required,max=255"`
}
type VideoCardData struct {
	Model    string `json:"model" validate:"required,max=255"`
	Producer string `json:"producer" validate:"required,max=255"`
}
type MonitorData struct {
	Model    string `json:"model" validate:"required,max=255"`
	Producer string `json:"producer" validate:"required,max=255"`
}
type RamData struct {
	Type     string `json:"type" validate:"required,max=255"`
	Capacity int    `json:"capacity" validate:"required,numeric"`
}
