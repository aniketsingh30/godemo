package qualysReport

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	SecurityScoreCardEntity "riscvue.com/pkg/entities/securityScoreCard"
	"riscvue.com/pkg/handlers"
	"riscvue.com/pkg/repository/adapter"
	Rules "riscvue.com/pkg/rules"
	SecurityScoreCardRules "riscvue.com/pkg/rules/securityScoreCard"
	"riscvue.com/pkg/services"
	QualysReportervice "riscvue.com/pkg/services/qualysReport"
	"riscvue.com/pkg/utils/logger"
	HttpStatus "riscvue.com/utils/http"
)

const CLASSS_NAME = "SecurityScoreCardHandler"

type Handler struct {
	handlers.Interface
	Repository adapter.Interface
	Req        *http.Request
	Rules      Rules.Interface
	Service    services.QualysReportInterface
}

func NewHandler(repository adapter.Interface, req *http.Request) handlers.Interface {
	return &Handler{
		Req:     req,
		Service: QualysReportervice.NewQualysReportService(repository),
		Rules:   SecurityScoreCardRules.NewRules(),
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	const MethodName = "Get"

	name := chi.URLParam(r, "domainName")
	name = "makenacap.com"
	all := chi.URLParam(r, "all")
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" request:::", name)
	response, err := h.Service.Get(name, all)
	if err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)
		return
	}

	HttpStatus.StatusOK(w, r, response)
}
func (h *Handler) CreateRecord(w http.ResponseWriter, r *http.Request) {
	productBody, err := h.getBodyAndValidate(r, uuid.Nil)
	if err != nil {
		HttpStatus.StatusBadRequest(w, r, err)
		return
	}
	h.Service.CreateRecord(productBody)
}

func (h *Handler) UpdateRecord(w http.ResponseWriter, r *http.Request) {
	productBody, err := h.getBodyAndValidate(r, uuid.Nil)
	if err != nil {
		HttpStatus.StatusBadRequest(w, r, err)
		return
	}
	h.Service.UpdateRecord(productBody)

}

func (h *Handler) DeleteRecord(w http.ResponseWriter, r *http.Request) {
	//return h.Service.DeleteRecord()
}
func (h *Handler) UnhandledMethod(w http.ResponseWriter, r *http.Request) {
	//return h.Service.UnhandledMethod()
}

func (h *Handler) getBodyAndValidate(r *http.Request, ID uuid.UUID) (*SecurityScoreCardEntity.MainResponse, error) {
	productBody := &SecurityScoreCardEntity.MainResponse{}
	body, err := h.Rules.ConvertIoReaderToStruct(r.Body, productBody)
	if err != nil {
		return &SecurityScoreCardEntity.MainResponse{}, errors.New("body is required")
	}

	productParsed, err := SecurityScoreCardEntity.InterfaceToModel(body)
	if err != nil {
		return &SecurityScoreCardEntity.MainResponse{}, errors.New("error on convert body to model")
	}

	setDefaultValues(productParsed, ID)

	return productParsed, h.Rules.Validate(productParsed)
}

func setDefaultValues(product *SecurityScoreCardEntity.MainResponse, ID uuid.UUID) {
	product.UpdatedAt = time.Now()
	if ID == uuid.Nil {
		product.ID = uuid.New()
		product.CreatedAt = time.Now()
	} else {
		product.ID = ID
	}
}
