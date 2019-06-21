package response

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/pinkgorilla/go-sample/starter-api/errors"
	"github.com/pinkgorilla/go-sample/starter-api/internal/alert"
)

// WithError emits proper error response
func WithError(w http.ResponseWriter, n alert.Alert, err error) {
	switch err.(type) {
	case errors.BaseError:
		JSON(w, http.StatusBadRequest, err)
	case errors.NotFoundError:
		JSON(w, http.StatusNotFound, err)
	case errors.CommonError:
		JSON(w, http.StatusBadRequest, err)
	case errors.ValidationError:
		JSON(w, http.StatusUnprocessableEntity, err)
	case errors.AuthError:
		JSON(w, http.StatusUnauthorized, err)
	case errors.PermissionError:
		JSON(w, http.StatusForbidden, err)
	case errors.ServiceError:
		logError(n, err)
		response := errors.NewBaseError(http.StatusText(http.StatusInternalServerError), "Server tidak dapat memproses permintaan anda, cobalah beberapa saat lagi.")
		JSON(w, http.StatusInternalServerError, response)
	default:
		logError(n, err)
		response := errors.NewBaseError(http.StatusText(http.StatusInternalServerError), "Server tidak dapat memproses permintaan anda, cobalah beberapa saat lagi.")
		JSON(w, http.StatusInternalServerError, response)
	}
}

func logError(a alert.Alert, err error) {
	msg := fmt.Sprintf("%+v\n%s", err, string(debug.Stack()))
	log.Println(msg)
	if a != nil {
		alert := alert.NewAlert(err.Error(), err, debug.Stack())
		if err := a.Alert(alert); err != nil {
			log.Println("Failed to alert using slack: ", err)
		}
	}
}
