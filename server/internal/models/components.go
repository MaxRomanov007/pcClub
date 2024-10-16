package models

// ProcessorProducer represents the processor_producers table.
type ProcessorProducer struct {
	ProcessorProducerID int64  `db:"processor_producer_id"`
	Name                string `db:"processor_producer"`
}

// VideoCardProducer represents the video_card_producers table.
type VideoCardProducer struct {
	VideoCardProducerID int64  `db:"video_card_producer_id"`
	Name                string `db:"video_card_producer"`
}

// MonitorProducer represents the monitor_producers table.
type MonitorProducer struct {
	MonitorProducerID int64  `db:"monitor_producer_id"`
	Name              string `db:"monitor_producer"`
}

// RamType represents the ram_types table.
type RamType struct {
	RamTypeID int64  `db:"ram_type_id"`
	Name      string `db:"ram_type"`
}

// Processor represents the processors table.
type Processor struct {
	ProcessorID         int64  `db:"processor_id"`
	ProcessorProducerID int64  `db:"processor_producer_id"`
	Model               string `db:"processor_model"`
	ProcessorProducer   *ProcessorProducer
}

// VideoCard represents the video_cards table.
type VideoCard struct {
	VideoCardID         int64  `db:"video_card_id"`
	VideoCardProducerID int64  `db:"video_card_producer_id"`
	Model               string `db:"video_card_model"`
	VideoCardProducer   *VideoCardProducer
}

// Monitor represents the monitors table.
type Monitor struct {
	MonitorID         int64  `db:"monitor_id"`
	MonitorProducerID int64  `db:"monitor_producer_id"`
	Model             string `db:"monitor_model"`
	MonitorProducer   *MonitorProducer
}

// Ram represents the ram table.
type Ram struct {
	RamID     int64 `db:"ram_id"`
	RamTypeID int64 `db:"ram_type_id"`
	Capacity  int   `db:"ram_capacity"`
	RamType   *RamType
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
