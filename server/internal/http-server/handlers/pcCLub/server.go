package pcCLub

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"server/internal/config"
	"server/internal/models"
)

type AuthService interface {
	Access(
		ctx context.Context,
		accessToken string,
	) (uid int64, err error)

	Refresh(
		ctx context.Context,
		refreshToken string,
	) (
		access string,
		refresh string,
		err error,
	)

	Tokens(
		ctx context.Context,
		uid int64,
	) (
		accessToken string,
		refreshToken string,
		err error,
	)

	BanTokens(
		ctx context.Context,
		accessToken string,
		refreshToken string,
	) (uid int64, err error)
}

type UserService interface {
	SaveUser(
		ctx context.Context,
		email string,
		password string,
	) (uid int64, err error)

	Login(
		ctx context.Context,
		email string,
		password string,
	) (uid int64, err error)

	User(
		ctx context.Context,
		uid int64,
	) (user models.User, err error)

	UserByEmail(
		ctx context.Context,
		email string,
	) (user models.User, err error)

	DeleteUser(
		ctx context.Context,
		uid int64,
	) (err error)

	IsAdmin(
		ctx context.Context,
		uid int64,
	) (err error)
}

type PcTypeService interface {
	PcTypes(
		ctx context.Context,
		limit int,
		offset int,
	) (pcs []models.PcType, err error)

	PcType(
		ctx context.Context,
		typeID int64,
	) (pcType models.PcType, err error)

	SavePcType(
		ctx context.Context,
		pc *models.PcType,
	) (id int64, err error)

	UpdatePcType(
		ctx context.Context,
		typeID int64,
		pcType *models.PcType,
	) (err error)

	DeletePcType(
		ctx context.Context,
		typeID int64,
	) (err error)
}

type PcService interface {
	Pcs(
		ctx context.Context,
		typeId int64,
		isAvailable bool,
	) (pcs []models.Pc, err error)

	SavePc(
		ctx context.Context,
		pc *models.Pc,
	) (id int64, err error)

	UpdatePc(
		ctx context.Context,
		pcID int64,
		pc *models.Pc,
	) (err error)

	DeletePc(
		ctx context.Context,
		pcId int64,
	) (err error)
}

type PcRoomService interface {
	PcRoom(
		ctx context.Context,
		roomID int64,
	) (room models.PcRoom, err error)

	PcRooms(
		ctx context.Context,
		pcTypeID int64,
	) (rooms []models.PcRoom, err error)

	SavePcRoom(
		ctx context.Context,
		room *models.PcRoom,
	) (id int64, err error)

	UpdatePcRoom(
		ctx context.Context,
		roomID int64,
		room *models.PcRoom,
	) (err error)

	DeletePcRoom(
		ctx context.Context,
		pcId int64,
	) (err error)
}

type ProcessorService interface {
	ProcessorProducers(
		ctx context.Context,
	) (producers []models.ProcessorProducer, err error)

	Processors(
		ctx context.Context,
		producerID int64,
	) (processors []models.Processor, err error)

	SaveProcessorProducer(
		ctx context.Context,
		producer *models.ProcessorProducer,
	) (id int64, err error)

	SaveProcessor(
		ctx context.Context,
		processor *models.Processor,
	) (id int64, err error)

	DeleteProcessorProducer(
		ctx context.Context,
		producerID int64,
	) (err error)

	DeleteProcessor(
		ctx context.Context,
		processorID int64,
	) (err error)
}

type MonitorService interface {
	MonitorProducers(
		ctx context.Context,
	) (producers []models.MonitorProducer, err error)

	Monitors(
		ctx context.Context,
		producerID int64,
	) (monitors []models.Monitor, err error)

	SaveMonitorProducer(
		ctx context.Context,
		producer *models.MonitorProducer,
	) (id int64, err error)

	SaveMonitor(
		ctx context.Context,
		monitor *models.Monitor,
	) (id int64, err error)

	DeleteMonitorProducer(
		ctx context.Context,
		producerID int64,
	) (err error)

	DeleteMonitor(
		ctx context.Context,
		monitorID int64,
	) (err error)
}

type VideoCardService interface {
	VideoCardProducers(
		ctx context.Context,
	) (producers []models.VideoCardProducer, err error)

	VideoCards(
		ctx context.Context,
		videoCardID int64,
	) (cards []models.VideoCard, err error)

	SaveVideoCardProducer(
		ctx context.Context,
		producer *models.VideoCardProducer,
	) (id int64, err error)

	SaveVideoCard(
		ctx context.Context,
		card *models.VideoCard,
	) (id int64, err error)

	DeleteVideoCardProducer(
		ctx context.Context,
		producerID int64,
	) (err error)

	DeleteVideoCard(
		ctx context.Context,
		cardID int64,
	) (err error)
}

type RamService interface {
	RamTypes(
		ctx context.Context,
	) (ramTypes []models.RAMType, err error)

	Rams(
		ctx context.Context,
		typeID int64,
	) (rams []models.RAM, err error)

	SaveRamType(
		ctx context.Context,
		ramType *models.RAMType,
	) (id int64, err error)

	SaveRam(
		ctx context.Context,
		ram *models.RAM,
	) (id int64, err error)

	DeleteRamType(
		ctx context.Context,
		ramTypeID int64,
	) (err error)

	DeleteRam(
		ctx context.Context,
		ramID int64,
	) (err error)
}

type ComponentsService struct {
	Processor ProcessorService
	Monitor   MonitorService
	VideoCard VideoCardService
	Ram       RamService
}

type DishService interface {
	Dishes(
		ctx context.Context,
		limit int,
		offset int,
	) (dishes []models.Dish, err error)

	Dish(
		ctx context.Context,
		dishId int64,
	) (dish models.Dish, err error)

	SaveDish(
		ctx context.Context,
		dish *models.Dish,
	) (id int64, err error)

	UpdateDish(
		ctx context.Context,
		dishID int64,
		dish *models.Dish,
	) (err error)

	DeleteDish(
		ctx context.Context,
		dishId int64,
	) (err error)
}

type OrderService interface {
	PcOrders(
		ctx context.Context,
		uid int64,
	)
}

type API struct {
	Log               *slog.Logger
	Cfg               *config.Config
	UserService       UserService
	AuthService       AuthService
	PcTypeService     PcTypeService
	PcService         PcService
	PcRoomService     PcRoomService
	ComponentsService ComponentsService
	DishService       DishService
	OrderService      OrderService
}

func New(
	log *slog.Logger,
	cfg *config.Config,
	userService UserService,
	authService AuthService,
	pcTypeService PcTypeService,
	pcService PcService,
	pcRoomService PcRoomService,
	componentsService ComponentsService,
	dishService DishService,
) *API {
	return &API{
		Log:               log,
		Cfg:               cfg,
		UserService:       userService,
		AuthService:       authService,
		PcTypeService:     pcTypeService,
		PcService:         pcService,
		PcRoomService:     pcRoomService,
		ComponentsService: componentsService,
		DishService:       dishService,
	}
}

func (a *API) log(op string, r *http.Request) *slog.Logger {
	return a.Log.With(
		slog.String("operation", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
}
