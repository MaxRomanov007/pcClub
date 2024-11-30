package pcCLub

import (
	"net/http"
	"time"
)

type OrderPcsRequest struct {
}

type SaveOrderPcRequest struct {
	PcId      int64     `json:"pc_id" validate:"required,min=1"`
	StartTime time.Time `json:"start_time" validate:"required,min=1"`
	Duration  int16     `json:"duration" validate:"required,min=1"`
}
type SaveOrderPcResponse struct {
	Code string `json:"code"`
}

func (a *API) OrderPcs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.pc-orderPc.OrderPc"
	}
}
