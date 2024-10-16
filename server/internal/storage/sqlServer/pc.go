package sqlServer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"server/internal/models"
)

type pcTypeFlat struct {
	TypeID            int64  `db:"type_id"`
	TypeName          string `db:"type"`
	ProcessorModel    string `db:"processor_model"`
	ProcessorProducer string `db:"processor_producer"`
	VideoCardModel    string `db:"video_card_model"`
	VideoCardProducer string `db:"video_card_producer"`
	MonitorModel      string `db:"monitor_model"`
	MonitorProducer   string `db:"monitor_producer"`
	RamType           string `db:"ram_type"`
	RamCapacity       int    `db:"ram_capacity"`
}

func (pcType pcTypeFlat) Parse() models.PcTypeData {
	return models.PcTypeData{
		TypeID:   pcType.TypeID,
		TypeName: pcType.TypeName,
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

func (s *Storage) Pcs(
	ctx context.Context,
	typeId int64,
	isAvailable bool,
) ([]models.PcData, error) {
	const op = "storage.sqlServer.pc.Pc"

	squirrelStmt := squirrel.Select(
		"pc.pc_id",
		"pc.row",
		"pc.place",
		"pc.description",
		"pc.pc_room_id",
	).
		From("pc").
		Where(squirrel.Eq{
			"pc.pc_type_id": typeId,
		})

	if isAvailable {
		squirrelStmt = squirrelStmt.
			Join("pc_statuses ON pc.pc_status_id = pc_statuses.pc_status_id").
			Where(squirrel.Eq{
				"pc_statuses.name": "available",
			})
	}

	stmt, args, err := squirrelStmt.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}

	stmt = replacePositionalParams(stmt, args)

	var pcs []models.PcData
	if err := s.db.SelectContext(ctx, &pcs, stmt, args...); err != nil {
		return nil, fmt.Errorf("%s: failed to get pcs: %w", op, handleError(err))
	}

	return pcs, nil
}

func (s *Storage) PcTypes(
	ctx context.Context,
	limit int64,
	offset int64,
) ([]models.PcTypeData, error) {
	const op = "storage.sqlServer.pc.pcTypes"

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
	const op = "storage.sqlServer.pc.SavePc"
	stmt, args, err := squirrel.
		Expr(
			"EXEC InsertPcType ?, ?, ?, ?, ?, ?, ?, ?, ?, ?",
			name,
			description,
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

func (s *Storage) SavePc(
	ctx context.Context,
	typeId int64,
	roomId int64,
	row int,
	place int,
) error {
	const op = "storage.sqlServer.pc.SavePc"

	stmt, args, err := squirrel.
		Insert("pc").
		Columns(
			"pc_type_id",
			"pc_room_id",
			"row",
			"place",
		).
		Values(
			typeId,
			roomId,
			row,
			place,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}

	stmt = replacePositionalParams(stmt, args)

	if _, err := s.db.ExecContext(ctx, stmt, args...); err != nil {
		return fmt.Errorf("%s: failed to exec statement: %w", op, handleError(err))
	}

	return nil
}
