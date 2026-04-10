package interfaces

import "net/http"

type HealthHandler interface {
	Handler
	HealthCheck(w http.ResponseWriter, r *http.Request)
}
