package orders

import (
	"net/http"
	"net/url"

	"github.com/Albert-Ti/go-diploma-tpl/internal/config"
)

func checkAccrualOrder(order string) int {
	urlAccrual := url.URL{
		Scheme: "http",
		Host:   config.Envs.AccrualSystemAddr,
		Path:   "api/orders/" + order,
	}
	res, _ := http.Get(urlAccrual.String())

	return res.StatusCode
}
