package orders

import (
	"io"
	"net/http"

	"github.com/Albert-Ti/go-diploma-tpl/internal/service"
	"github.com/Albert-Ti/go-diploma-tpl/internal/utils"
)

func AddOrders(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		b, _ := io.ReadAll(r.Body)

		if !utils.AlgoLuna(string(b)) {
			http.Error(w, "Incorrect order number format", http.StatusUnprocessableEntity)
			return
		}
	}
}
