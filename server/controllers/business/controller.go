package business

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pinkgorilla/go-sample/internal/business"
	"github.com/pinkgorilla/go-sample/internal/event"
)

// Controller ...
type Controller struct {
	Router   *chi.Mux
	business *business.Service
	log      func(e event.Event)
}

// NewController ...
func NewController(business *business.Service, log func(e event.Event)) *Controller {

	controller := &Controller{
		Router:   chi.NewRouter(),
		business: business,
		log:      log,
	}

	controller.Router.Get("/token", controller.GetToken())
	return controller
}

// GetToken ...
func (controller *Controller) GetToken() func(writer http.ResponseWriter, request *http.Request) {
	e := event.Event{
		Name:   "HANDLER_CREATED",
		Source: "controller.GetToken",
		Data:   controller,
	}
	controller.log(e)

	// GetToken , handler function for getting a token
	return func(writer http.ResponseWriter, request *http.Request) {
		e := event.Event{
			Name:   "REQUEST_HANDLER_EXECUTED",
			Source: "controller.GetToken",
			Data:   controller,
		}
		controller.log(e)

		username, password := "", ""
		actor := "go-sample"
		service := controller.business
		service.Actor = actor
		token := service.GetToken(username, password)
		writer.Write([]byte(token))
	}
}
