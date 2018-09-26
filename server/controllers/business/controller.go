package business

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pinkgorilla/go-sample/internal/business"
)

// Controller ...
type Controller struct {
	Router   *chi.Mux
	business *business.Service
}

// NewController ...
func NewController(business *business.Service) *Controller {

	controller := &Controller{
		Router:   chi.NewRouter(),
		business: business,
	}

	controller.Router.Get("/token", controller.GetToken())
	return controller
}

// GetToken ...
func (controller *Controller) GetToken() func(writer http.ResponseWriter, request *http.Request) {
	// GetToken , handler function for getting a token
	return func(writer http.ResponseWriter, request *http.Request) {
		username, password := "", ""
		service := controller.business
		token := service.GetToken(username, password)
		writer.Write([]byte(token))
	}
}
