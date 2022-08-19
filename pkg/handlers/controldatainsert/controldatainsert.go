package controldatainsert

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	EntityControlData "riscvue.com/pkg/entities/controldata"
	"riscvue.com/pkg/handlers"
	"riscvue.com/pkg/repository/adapter"
	Rules "riscvue.com/pkg/rules"
	RulesControldata "riscvue.com/pkg/rules/controldata"
	"riscvue.com/pkg/services"
	ControldataInsertService "riscvue.com/pkg/services/controldatainsert"
	HttpStatus "riscvue.com/utils/http"
)

const CLASSS_NAME = "Handler"

type Handler struct {
	handlers.ControlDataInsertInterface
	Repository adapter.Interface
	Req        *http.Request
	Rules      Rules.Interface
	Service    services.ControlDataInsertInterface
}

func NewHandler(repository adapter.Interface, req *http.Request) handlers.ControlDataInsertInterface {
	return &Handler{
		Req:     req,
		Service: ControldataInsertService.NewControlDataInsertService(repository),
		Rules:   RulesControldata.NewRules(),
	}
}

func (h *Handler) CreateRecord(w http.ResponseWriter, r *http.Request) {
	productBody, err := h.getBodyAndValidate(r, uuid.Nil)
	if err != nil {
		HttpStatus.StatusBadRequest(w, r, err)
		return
	}
	h.Service.CreateRecord(productBody)
}

func (h *Handler) CreateRecordFromOtherEnv(w http.ResponseWriter, r *http.Request) {
	env := chi.URLParam(r, "env")
	datas, err := h.Service.CreateRecordFromOtherEnv(env)
	if err != nil {
		HttpStatus.StatusBadRequest(w, r, err)
		return
	}

	HttpStatus.StatusOK(w, r, datas)

}

func (h *Handler) getBodyAndValidate(r *http.Request, ID uuid.UUID) (*EntityControlData.ControlDataRequest, error) {
	productBody := &EntityControlData.ControlDataRequest{}
	body, err := h.Rules.ConvertIoReaderToStruct(r.Body, productBody)
	if err != nil {
		return &EntityControlData.ControlDataRequest{}, errors.New("body is required")
	}

	productParsed, err := EntityControlData.InterfaceToModelReq(body)
	if err != nil {
		return &EntityControlData.ControlDataRequest{}, errors.New("error on convert body to model")
	}

	return productParsed, h.Rules.Validate(productParsed)
}
