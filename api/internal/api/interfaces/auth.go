package interfaces

import "net/http"

type AuthHandler interface {
	Handler
	Me(w http.ResponseWriter, r *http.Request)
}
