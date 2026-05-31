package orders

import (
	"fmt"
	"net/http"

	"github.com/Albert-Ti/go-diploma-tpl/internal/service"
)

func AddOrders(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("ADD ORDERS", r.Context().Value("userID"))
	}
}
