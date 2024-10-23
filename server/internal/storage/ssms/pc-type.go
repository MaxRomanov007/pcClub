package ssms

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"server/internal/models"
)

type pcTypeFlat struct {
	TypeID            int64            `db:"type_id"`
	TypeName          string           `db:"type"`
	Description       sql.Null[string] `db:"description"`
	ProcessorModel    string           `db:"processor_model"`
	ProcessorProducer string           `db:"processor_producer"`
	VideoCardModel    string           `db:"video_card_model"`
	VideoCardProducer string           `db:"video_card_producer"`
	MonitorModel      string           `db:"monitor_model"`
	MonitorProducer   string           `db:"monitor_producer"`
	RamType           string           `db:"ram_type"`
	RamCapacity       int              `db:"ram_capacity"`
}

func (pcType pcTypeFlat) Parse() models.PcTypeData {
	return models.PcTypeData{
		TypeID:      pcType.TypeID,
		TypeName:    pcType.TypeName,
		Description: pcType.Description.V,
		Processor: &models.ProcessorData{
			Model:    pcType.ProcessorModel,
			Producer: pcType.ProcessorProducer,
		},
		VideoCard: &models.VideoCardData{
			Model:    pcType.VideoCardModel,
			Producer: pcType.VideoCardProducer,
		},
		Monitor: &models.MonitorData{
			Model:    pcType.MonitorModel,
			Producer: pcType.MonitorProducer,
		},
		Ram: &models.RamData{
			Capacity: pcType.RamCapacity,
			Type:     pcType.RamType,
		},
	}
}

func squirrelSelectPcTypeFlat() squirrel.SelectBuilder {
	return squirrel.Select(
		"pc_types.pc_type_id AS type_id",
		"pc_types.name AS type",
		"pc_types.description AS description",
		"processors.model AS processor_model",
		"processor_producers.name AS processor_producer",
		"video_cards.model AS video_card_model",
		"video_card_producers.name AS video_card_producer",
		"monitors.model AS monitor_model",
		"monitor_producers.name AS monitor_producer",
		"ram.capacity AS ram_capacity",
		"ram_types.name AS ram_type",
	).
		From("pc_types").
		Join("processors ON pc_types.processor_id = processors.processor_id").
		Join("video_cards ON pc_types.video_card_id = video_cards.video_card_id").
		Join("monitors ON pc_types.monitor_id = monitors.monitor_id").
		Join("ram ON pc_types.ram_id = ram.ram_id").
		Join("processor_producers ON processors.processor_producer_id = processor_producers.processor_producer_id").
		Join("video_card_producers ON video_cards.video_card_producer_id = video_card_producers.video_card_producer_id").
		Join("monitor_producers ON monitors.monitor_producer_id = monitor_producers.monitor_producer_id").
		Join("ram_types ON ram.ram_type_id = ram_types.ram_type_id")
}

func (s *Storage) PcTypes(
	ctx context.Context,
	limit int64,
	offset int64,
) ([]models.PcTypeData, error) {
	const op = "storage.ssms.pc.pcTypes"

	squirrelStmt := squirrelSelectPcTypeFlat()

	if limit > 0 || offset > 0 {
		squirrelStmt = squirrelStmt.
			Column("ROW_NUMBER() OVER (ORDER BY pc_types.pc_type_id) AS nom").
			Prefix("WITH num_row AS (").
			Suffix(")")

		stmt, args, err := squirrelStmt.ToSql()
		if err != nil {
			return nil, fmt.Errorf("%s: failed to prepare inner statement: %w", op, err)
		}

		squirrelStmt = squirrel.Select(
			"type_id",
			"type",
			"description",
			"processor_model",
			"processor_producer",
			"monitor_model",
			"video_card_model",
			"video_card_producer",
			"monitor_model",
			"monitor_producer",
			"ram_capacity",
			"ram_type",
		).
			From("num_row").
			Where("nom BETWEEN ? AND ?", offset+1, limit+offset).
			Prefix(stmt, args...)
	}

	stmt, args, err := squirrelStmt.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}

	stmt = replacePositionalParams(stmt, args)

	var pcTypesFlat []pcTypeFlat
	if err := s.db.SelectContext(ctx, &pcTypesFlat, stmt, args...); err != nil {
		return nil, fmt.Errorf("%s: failed to get pc types: %w", op, handleError(err))
	}

	pcTypes := make([]models.PcTypeData, len(pcTypesFlat))
	for i, pcType := range pcTypesFlat {
		pcTypes[i] = pcType.Parse()
	}

	return pcTypes, nil
}

func (s *Storage) PcType(
	ctx context.Context,
	id int64,
) (models.PcTypeData, error) {
	const op = "storage.pcType.pcType"

	squirrelStmt := squirrelSelectPcTypeFlat()
	squirrelStmt = squirrelStmt.Where(squirrel.Eq{
		"pc_types.pc_type_id": id,
	})

	stmt, args, err := squirrelStmt.ToSql()
	if err != nil {
		return models.PcTypeData{}, fmt.Errorf("%s: failed to prepare statement: %w", op, handleError(err))
	}

	stmt = replacePositionalParams(stmt, args)

	var pcType pcTypeFlat
	if err := s.db.GetContext(ctx, &pcType, stmt, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.PcTypeData{}, fmt.Errorf("%s: failed to get pc type: %w", op, ErrNotFound)
		}
		return models.PcTypeData{}, fmt.Errorf("%s: failed to get pc type: %w", op, handleError(err))
	}

	return pcType.Parse(), nil
}

func (s *Storage) SavePcType(
	ctx context.Context,
	name string,
	description string,
	processor *models.ProcessorData,
	videoCard *models.VideoCardData,
	monitor *models.MonitorData,
	ram *models.RamData,
) error {
	const op = "storage.ssms.pc.SavePc"

	stmt, args, err := squirrel.
		Expr(
			"EXEC InsertPcType ?, ?, ?, ?, ?, ?, ?, ?, ?, ?",
			name,
			nullString(description),
			processor.Producer,
			processor.Model,
			videoCard.Producer,
			videoCard.Model,
			monitor.Producer,
			monitor.Model,
			ram.Type,
			ram.Capacity,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to prepare statement: %w", op, handleError(err))
	}

	stmt = replacePositionalParams(stmt, args)

	if _, err := s.db.ExecContext(ctx, stmt, args...); err != nil {
		return fmt.Errorf("%s: failed to exec statement: %w", op, handleError(err))
	}

	return nil
}

func (s *Storage) UpdatePcType(
	ctx context.Context,
	typeId int64,
	name string,
	description string,
	processor *models.ProcessorData,
	videoCard *models.VideoCardData,
	monitor *models.MonitorData,
	ram *models.RamData,
) error {
	const op = "storage.ssms.pc.UpdatePc"

	stmt, args, err := squirrel.Expr(
		"EXEC UpdatePcType ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?",
		typeId,
		name,
		nullString(description),
		processor.Producer,
		processor.Model,
		videoCard.Producer,
		videoCard.Model,
		monitor.Producer,
		monitor.Model,
		ram.Type,
		ram.Capacity,
	).ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}

	stmt = replacePositionalParams(stmt, args)

	var affectedCount int64
	if err := s.db.GetContext(ctx, &affectedCount, stmt, args...); err != nil {
		return fmt.Errorf("%s: failed to update pc type: %w", op, handleError(err))
	}

	if affectedCount == 0 {
		return fmt.Errorf("%s: %w", op, ErrNotFound)
	}

	return nil
}

func (s *Storage) DeletePcType(
	ctx context.Context,
	typeId int64,
) error {
	const op = "storage.ssms.pc.DeletePc"

	stmt, args, err := squirrel.
		Delete("pc_types").
		Where(squirrel.Eq{
			"pc_type_id": typeId,
		}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to prepare statement: %w", op, handleError(err))
	}

	stmt = replacePositionalParams(stmt, args)

	if _, err := s.db.ExecContext(ctx, stmt, args...); err != nil {
		return fmt.Errorf("%s: failed to exec statement: %w", op, handleError(err))
	}

	return nil
}
