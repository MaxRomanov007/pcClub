package main

import (
	"gorm.io/driver/sqlserver"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"log"
	"server/internal/config"
	"server/internal/lib/api/database/mssql"
)

func main() {
	cfg := config.MustLoad()
	db, err := gorm.Open(sqlserver.Open(mssql.GenerateConnString(cfg.Database.SQLServer)), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: " + err.Error())
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:      "./internal/models",
		ModelPkgPath: "./models",
		Mode:         gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.UseDB(db)

	var dataMap = map[string]func(gorm.ColumnType) string{
		"DECIMAL": func(dataType gorm.ColumnType) string {
			return "float32"
		},
		"INT": func(dataType gorm.ColumnType) string {
			return "int"
		},
		"SMALLINT": func(dataType gorm.ColumnType) string {
			return "int16"
		},
	}
	g.WithDataTypeMap(dataMap)

	g.GenerateAllTable()

	userRoles := g.GenerateModel("user_roles",
		gen.FieldRelate(field.HasMany, "Users", g.GenerateModel("users"), &field.RelateConfig{}),
	)
	users := g.GenerateModel("users",
		gen.FieldRelate(field.BelongsTo, "UserRole", userRoles, &field.RelateConfig{}),
		gen.FieldRelate(field.HasMany, "PcOrders", g.GenerateModel("pc_orders"), &field.RelateConfig{}),
		gen.FieldRelate(field.HasMany, "DishOrders", g.GenerateModel("dish_orders"), &field.RelateConfig{}),
	)

	processorProducers := g.GenerateModel("processor_producers",
		gen.FieldRelate(field.HasMany, "Processors", g.GenerateModel("processors"), &field.RelateConfig{}),
	)
	processors := g.GenerateModel("processors",
		gen.FieldRelate(field.BelongsTo, "ProcessorProducer", processorProducers, &field.RelateConfig{}),
		gen.FieldRelate(field.HasMany, "PcTypes", g.GenerateModel("pc_types"), &field.RelateConfig{}),
	)

	videoCardProducers := g.GenerateModel("video_card_producers",
		gen.FieldRelate(field.HasMany, "VideoCards", g.GenerateModel("video_cards"), &field.RelateConfig{}),
	)
	videoCards := g.GenerateModel("video_cards",
		gen.FieldRelate(field.BelongsTo, "VideoCardProducer", videoCardProducers, &field.RelateConfig{}),
		gen.FieldRelate(field.HasMany, "PcTypes", g.GenerateModel("pc_types"), &field.RelateConfig{}),
	)

	monitorProducers := g.GenerateModel("monitor_producers",
		gen.FieldRelate(field.HasMany, "Monitors", g.GenerateModel("monitors"), &field.RelateConfig{}),
	)
	monitors := g.GenerateModel("monitors",
		gen.FieldRelate(field.BelongsTo, "MonitorProducer", monitorProducers, &field.RelateConfig{}),
		gen.FieldRelate(field.HasMany, "PcTypes", g.GenerateModel("pc_types"), &field.RelateConfig{}),
	)

	ramTypes := g.GenerateModel("ram_types",
		gen.FieldRelate(field.HasMany, "Ram", g.GenerateModel("ram"), &field.RelateConfig{}),
	)
	ram := g.GenerateModel("ram",
		gen.FieldRelate(field.BelongsTo, "RAMType", ramTypes, &field.RelateConfig{}),
		gen.FieldRelate(field.HasMany, "PcTypes", g.GenerateModel("pc_types"), &field.RelateConfig{}),
	)

	pcTypes := g.GenerateModel("pc_types",
		gen.FieldRelate(field.BelongsTo, "Processor", processors, &field.RelateConfig{}),
		gen.FieldRelate(field.BelongsTo, "VideoCard", videoCards, &field.RelateConfig{}),
		gen.FieldRelate(field.BelongsTo, "Monitor", monitors, &field.RelateConfig{}),
		gen.FieldRelate(field.BelongsTo, "RAM", ram, &field.RelateConfig{}),
		gen.FieldRelate(field.HasMany, "Pcs", g.GenerateModel("pc"), &field.RelateConfig{}),
		gen.FieldRelate(field.HasMany, "PcTypeImages", g.GenerateModel("pc_type_images"), &field.RelateConfig{}),
	)

	g.GenerateModel("pc_type_images",
		gen.FieldRelate(field.BelongsTo, "PcType", pcTypes, &field.RelateConfig{}),
	)

	pcRooms := g.GenerateModel("pc_rooms",
		gen.FieldRelate(field.HasMany, "Pcs", g.GenerateModel("pc"), &field.RelateConfig{}),
	)
	pcStatuses := g.GenerateModel("pc_statuses",
		gen.FieldRelate(field.HasMany, "Pcs", g.GenerateModel("pc"), &field.RelateConfig{}),
	)
	pc := g.GenerateModel("pc",
		gen.FieldRelate(field.BelongsTo, "PcRoom", pcRooms, &field.RelateConfig{}),
		gen.FieldRelate(field.BelongsTo, "PcType", pcTypes, &field.RelateConfig{}),
		gen.FieldRelate(field.BelongsTo, "PcStatus", pcStatuses, &field.RelateConfig{}),
		gen.FieldRelate(field.HasMany, "PcOrders", g.GenerateModel("pc_orders"), &field.RelateConfig{}),
	)

	pcOrderStatuses := g.GenerateModel("pc_order_statuses",
		gen.FieldRelate(field.HasMany, "PcOrders", g.GenerateModel("pc_orders"), &field.RelateConfig{}),
	)
	g.GenerateModel("pc_orders",
		gen.FieldRelate(field.BelongsTo, "User", users, &field.RelateConfig{}),
		gen.FieldRelate(field.BelongsTo, "Pc", pc, &field.RelateConfig{}),
		gen.FieldRelate(field.BelongsTo, "PcOrderStatus", pcOrderStatuses, &field.RelateConfig{}),
	)

	dishStatuses := g.GenerateModel("dish_statuses",
		gen.FieldRelate(field.HasMany, "Dishes", g.GenerateModel("dishes"), &field.RelateConfig{}),
	)
	dishes := g.GenerateModel("dishes",
		gen.FieldRelate(field.BelongsTo, "DishStatus", dishStatuses, &field.RelateConfig{}),
		gen.FieldRelate(field.HasMany, "DishImages", g.GenerateModel("dish_images"), &field.RelateConfig{}),
		gen.FieldRelate(field.HasMany, "DishOrderList", g.GenerateModel("dish_order_list"), &field.RelateConfig{}),
	)

	g.GenerateModel("dish_images",
		gen.FieldRelate(field.BelongsTo, "Dish", dishes, &field.RelateConfig{}),
	)

	dishOrderStatuses := g.GenerateModel("dish_order_statuses",
		gen.FieldRelate(field.HasMany, "DishOrders", g.GenerateModel("dish_orders"), &field.RelateConfig{}),
	)
	dishOrders := g.GenerateModel("dish_orders",
		gen.FieldRelate(field.BelongsTo, "DishOrderStatus", dishOrderStatuses, &field.RelateConfig{}),
		gen.FieldRelate(field.BelongsTo, "User", users, &field.RelateConfig{}),
		gen.FieldRelate(field.HasMany, "DishOrderList", g.GenerateModel("dish_order_list"), &field.RelateConfig{}),
	)

	g.GenerateModel("dish_order_list",
		gen.FieldRelate(field.BelongsTo, "DishOrder", dishOrders, &field.RelateConfig{}),
		gen.FieldRelate(field.BelongsTo, "Dish", dishes, &field.RelateConfig{}),
	)

	g.Execute()
}
