package assessment

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	EntityAssessment "riscvue.com/pkg/entities/assessment"
	"riscvue.com/pkg/handlers"
	"riscvue.com/pkg/repository/adapter"
	Rules "riscvue.com/pkg/rules"
	RulesAssessment "riscvue.com/pkg/rules/assessment"
	"riscvue.com/pkg/services"
	AssessmentService "riscvue.com/pkg/services/assessment"
	"riscvue.com/pkg/utils/logger"
	HttpStatus "riscvue.com/utils/http"
)

const CLASSS_NAME = "Handler"

type Handler struct {
	handlers.Interface
	Repository adapter.Interface
	Req        *http.Request
	Rules      Rules.Interface
	Service    services.AssessmentInterface
}

type EntriesResponse struct {
	Entries []Entries `json:"entries"`
	Total   int64     `json:"total"`
}
type Entries struct {
	Name         string         `json:"name"`
	Score        int64          `json:"score"`
	Grade        string         `json:"grade"`
	Grade_url    string         `json:"grade_url"`
	IssueSummary []IssueSummary `json:"issue_summary"`
}
type IssueSummary struct {
	Type               string `json:"type"`
	Count              int64  `json:"count"`
	Severity           string `json:"severity"`
	Total_Score_Impact string `json:"total_score_impact"`
	Detail_Url         string `json:"detail_url"`
}
type Response struct {
	Name                   string `json:"name"`
	Description            string `json:"description"`
	Domain                 string `json:"domain"`
	Grade_url              string `json:"grade_url"`
	Industry               string `json:"industry"`
	Size                   string `json:"size"`
	Score                  int64  `json:"score"`
	Grade                  string `json:"grade"`
	Last30day_score_change int64  `json:"last30day_score_change"`
	Created_at             string `json:"created_at"`
	Disputed               bool   `json:"disputed"`
}

func NewHandler(repository adapter.Interface, req *http.Request) handlers.Interface {
	return &Handler{
		Req:     req,
		Service: AssessmentService.NewAssessmentService(repository),
		Rules:   RulesAssessment.NewRules(),
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	const MethodName = "Get"

	client1 := &http.Client{}
	req1, err := http.NewRequest("GET", "https://platform-api.securityscorecard.io/companies/makenacap.com/", nil)
	if err != nil {
		fmt.Print(err.Error())
	}
	req1.Header.Add("Accept", "application/json")
	req1.Header.Add("Content-Type", "application/json")
	req1.Header.Add("Authorization", "Token ZJfQ888QHGnCkZySsLqkOp1DSklH")
	resp1, err := client1.Do(req1)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp1.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp1.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	var responseObject Response
	json.Unmarshal(bodyBytes, &responseObject)
	fmt.Printf("API Response as struct %+v\n", responseObject)

	//client := &http.Client{}
	req, err := http.NewRequest("GET", "https://platform-api.securityscorecard.io/companies/makenacap.com/factors", nil)
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Token ZJfQ888QHGnCkZySsLqkOp1DSklH")
	resp, err := client1.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()
	bodyBytes1, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	var EntriesResponse EntriesResponse
	json.Unmarshal(bodyBytes1, &EntriesResponse)
	fmt.Printf("API EntriesResponse as struct %+v\n", EntriesResponse)

	ID := chi.URLParam(r, "ID")

	history := chi.URLParam(r, "includeHistory")
	userId := chi.URLParam(r, "userId")
	createdBy := chi.URLParam(r, "createdBy")
	name := chi.URLParam(r, "name")
	customerId := chi.URLParam(r, "customerId")
	all := chi.URLParam(r, "all")
	logger.INFO("className="+CLASSS_NAME+" MethodName="+MethodName+" request:::", ID)
	response, err := h.Service.Get(ID, history, createdBy, name, customerId, userId, all)
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

func (h *Handler) getBodyAndValidate(r *http.Request, ID uuid.UUID) (*EntityAssessment.Assessment, error) {
	productBody := &EntityAssessment.Assessment{}
	body, err := h.Rules.ConvertIoReaderToStruct(r.Body, productBody)
	if err != nil {
		return &EntityAssessment.Assessment{}, errors.New("body is required")
	}

	productParsed, err := EntityAssessment.InterfaceToModel(body)
	if err != nil {
		return &EntityAssessment.Assessment{}, errors.New("error on convert body to model")
	}

	setDefaultValues(productParsed, ID)

	return productParsed, h.Rules.Validate(productParsed)
}

func setDefaultValues(product *EntityAssessment.Assessment, ID uuid.UUID) {
	product.UpdatedAt = time.Now()
	if ID == uuid.Nil {
		product.ID = uuid.New()
		product.CreatedAt = time.Now()
	} else {
		product.ID = ID
	}
}
