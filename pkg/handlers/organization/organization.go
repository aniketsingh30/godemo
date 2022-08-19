package organization

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	EntityCustomer "riscvue.com/pkg/entities/customer"
	"riscvue.com/pkg/handlers"
	"riscvue.com/pkg/repository/adapter"
	Rules "riscvue.com/pkg/rules"
	OrganizationRules "riscvue.com/pkg/rules/organization"
	"riscvue.com/pkg/services"
	OrganizationService "riscvue.com/pkg/services/organization"
	HttpStatus "riscvue.com/utils/http"
)

const CLASSS_NAME = "Handler"

type Handler struct {
	handlers.Interface
	Repository adapter.Interface
	Req        *http.Request
	Service    services.OrganizationInterface
	Rules      Rules.Interface
}

func NewHandler(repository adapter.Interface, req *http.Request) handlers.Interface {
	return &Handler{
		Req:     req,
		Service: OrganizationService.NewOrganizationService(repository),
		Rules:   OrganizationRules.NewRules(),
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	response, err := h.Service.GetRecordByCustomerId("690659aa-9138-4b3d-9b75-ba8e5bed9b95")
	if err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)

	}
	HttpStatus.StatusOK(w, r, response)
}
func (h *Handler) CreateRecord(w http.ResponseWriter, r *http.Request) {
	productBody, err := h.getBodyAndValidate(r, uuid.Nil)
	if err != nil {
		HttpStatus.StatusBadRequest(w, r, err)
		return
	}
	response, err := h.Service.CreateRecord(productBody)
	if err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)

	}
	HttpStatus.StatusOK(w, r, response)
}

func (h *Handler) UpdateRecord(w http.ResponseWriter, r *http.Request) {
	productBody, err := h.getBodyAndValidateUpdate(r, uuid.Nil)
	if err != nil {
		HttpStatus.StatusBadRequest(w, r, err)
		return
	}
	response, err := h.Service.UpdateRecord(productBody)
	if err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)

	}
	HttpStatus.StatusOK(w, r, response)

}

func (h *Handler) DeleteRecord(w http.ResponseWriter, r *http.Request) {
	response, err := h.Service.DeleteRecord("f")
	if err != nil {
		HttpStatus.StatusInternalServerError(w, r, err)

	}
	HttpStatus.StatusOK(w, r, response)
}
func (h *Handler) getBodyAndValidate(r *http.Request, ID uuid.UUID) (*EntityCustomer.Customer, error) {
	productBody := &EntityCustomer.Customer{}
	body, err := h.Rules.ConvertIoReaderToStruct(r.Body, productBody)
	if err != nil {
		return &EntityCustomer.Customer{}, errors.New("body is required")
	}

	productParsed, err := EntityCustomer.InterfaceToModel(body)
	if err != nil {
		return &EntityCustomer.Customer{}, errors.New("error on convert body to model")
	}

	setDefaultValues(productParsed, ID)

	return productParsed, h.Rules.Validate(productParsed)
}

func (h *Handler) getBodyAndValidateUpdate(r *http.Request, ID uuid.UUID) (*EntityCustomer.CustomerUpdate, error) {
	productBody := &EntityCustomer.CustomerUpdate{}
	body, err := h.Rules.ConvertIoReaderToStruct(r.Body, productBody)
	if err != nil {
		return &EntityCustomer.CustomerUpdate{}, errors.New("body is required")
	}

	productParsed, err := EntityCustomer.InterfaceToModelUpdate(body)
	if err != nil {
		return &EntityCustomer.CustomerUpdate{}, errors.New("error on convert body to model")
	}

	return productParsed, h.Rules.ValidateUpdate(productParsed)
}
func setDefaultValues(assesment *EntityCustomer.Customer, ID uuid.UUID) {
	assesment.UpdatedAt = time.Now()
	if ID == uuid.Nil {
		assesment.CustomerId = uuid.New().String()
		assesment.CreatedAt = time.Now()
	} else {
		assesment.CustomerId = ID.String()
	}

}
