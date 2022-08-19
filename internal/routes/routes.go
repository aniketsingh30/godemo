package routes

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	ServerConfig "riscvue.com/config"
	HealthHandler "riscvue.com/internal/handlers/health"
	ProductHandler "riscvue.com/internal/handlers/product"
	AssessmentHandler "riscvue.com/pkg/handlers/assessment"
	ControldatainsertHandler "riscvue.com/pkg/handlers/controldatainsert"
	OrganizationHandler "riscvue.com/pkg/handlers/organization"
	QualysReportHandler "riscvue.com/pkg/handlers/qualysReport"
	SecurityScoreCardHandler "riscvue.com/pkg/handlers/securityScoreCard"
	"riscvue.com/pkg/repository/adapter"
)

type Router struct {
	config *Config
	router *chi.Mux
}

func NewRouter() *Router {
	return &Router{
		config: NewConfig().SetTimeout(ServerConfig.GetConfig().Timeout),
		router: chi.NewRouter(),
	}
}

func (r *Router) SetRouters(repository adapter.Interface) *chi.Mux {
	r.setConfigsRouters()

	r.RouterHealth(repository)
	r.RouterProduct(repository)
	r.RouterAssessment(repository)
	r.RouterControlData(repository)
	r.RouterSecurity(repository)
	r.RouterQualys(repository)
	r.RouterOrganization(repository)

	return r.router
}

func (r *Router) setConfigsRouters() {
	r.EnableCORS()
	r.EnableLogger()
	r.EnableTimeout()
	r.EnableRecover()
	r.EnableRequestID()
	r.EnableRealIP()
}

func (r *Router) RouterHealth(repository adapter.Interface) {
	handler := HealthHandler.NewHandler(repository)

	r.router.Route("/health", func(route chi.Router) {
		route.Post("/", handler.Post)
		route.Get("/", handler.Get)
		route.Put("/", handler.Put)
		route.Delete("/", handler.Delete)
		route.Options("/", handler.Options)
	})
}

func (r *Router) RouterProduct(repository adapter.Interface) {
	handler := ProductHandler.NewHandler(repository)

	r.router.Route("/product", func(route chi.Router) {
		route.Post("/", handler.Post)
		route.Get("/", handler.Get)
		route.Get("/{ID}", handler.Get)
		route.Put("/{ID}", handler.Put)
		route.Delete("/{ID}", handler.Delete)
		route.Options("/", handler.Options)
	})
}

func (r *Router) RouterAssessment(repository adapter.Interface) {
	handler := AssessmentHandler.NewHandler(repository, &http.Request{})

	r.router.Route("/assessment", func(route chi.Router) {
		route.Post("/", handler.CreateRecord)
		route.Get("/", handler.Get)

		route.Put("/", handler.UpdateRecord)
		route.Delete("/{ID}", handler.DeleteRecord)

	})
}

func (r *Router) RouterSecurity(repository adapter.Interface) {
	handler := SecurityScoreCardHandler.NewHandler(repository, &http.Request{})

	r.router.Route("/security", func(route chi.Router) {
		route.Post("/", handler.CreateRecord)
		route.Get("/", handler.Get)

		route.Put("/", handler.UpdateRecord)
		route.Delete("/{ID}", handler.DeleteRecord)

	})
}

func (r *Router) RouterQualys(repository adapter.Interface) {
	handler := QualysReportHandler.NewHandler(repository, &http.Request{})

	r.router.Route("/qualys", func(route chi.Router) {
		route.Post("/", handler.CreateRecord)
		route.Get("/", handler.Get)

		route.Put("/", handler.UpdateRecord)
		route.Delete("/{ID}", handler.DeleteRecord)

	})
}

func (r *Router) RouterOrganization(repository adapter.Interface) {
	handler := OrganizationHandler.NewHandler(repository, &http.Request{})

	r.router.Route("/customer", func(route chi.Router) {
		route.Post("/", handler.CreateRecord)
		route.Get("/", handler.Get)

		route.Put("/", handler.UpdateRecord)
		route.Delete("/{ID}", handler.DeleteRecord)

	})
}
func (r *Router) RouterControlData(repository adapter.Interface) {
	handler := ControldatainsertHandler.NewHandler(repository, &http.Request{})

	r.router.Route("/controlData", func(route chi.Router) {
		route.Post("/", handler.CreateRecord)
		route.Get("/", handler.CreateRecordFromOtherEnv)

	})
}

func (r *Router) EnableLogger() *Router {
	r.router.Use(middleware.Logger)
	return r
}

func (r *Router) EnableTimeout() *Router {
	r.router.Use(middleware.Timeout(r.config.GetTimeout()))
	return r
}

func (r *Router) EnableCORS() *Router {
	r.router.Use(r.config.Cors)
	return r
}

func (r *Router) EnableRecover() *Router {
	r.router.Use(middleware.Recoverer)
	return r
}

func (r *Router) EnableRequestID() *Router {
	r.router.Use(middleware.RequestID)
	return r
}

func (r *Router) EnableRealIP() *Router {
	r.router.Use(middleware.RealIP)
	return r
}
