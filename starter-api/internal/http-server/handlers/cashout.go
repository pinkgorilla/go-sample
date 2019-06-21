package handlers

import (
	"net/http"
	"time"

	"github.com/pinkgorilla/go-sample/starter-api/internal/cashout"
)

// CreateCashOutHandler swagger:operation POST /v1/cashout handlers CreateCashOutHandler
//
// return user acl
//
// ---
// produces:
// - application/json
// responses:
//   200:
//     description: "Ok"
//     schema:
//       $ref: "#/definitions/CreateCashOutResponse"
//
func CreateCashOutHandler() http.HandlerFunc {

	// CreateCashOutResponse ...
	// swagger:model
	type payload struct {
		TransactionCode string          `json:"trxcode"`
		Sender          *cashout.Person `json:"sender"`
		Recipient       *cashout.Person `json:"recipient"`
		Currency        string          `json:"currency"`
		Amount          float64         `json:"amount"`
		IssuedDate      time.Time       `json:"issuedDate"`
		ExpiryDate      time.Time       `json:"expiryDate"`
	}

	return Handler(func(w http.ResponseWriter, r *http.Request) error {
		// return errors.NewCommonError("common-error")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world"))
		return nil
	})
}
