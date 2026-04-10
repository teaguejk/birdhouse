package responses

import (
	"bufio"
	"encoding/json"
	"net"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func TokenError(w http.ResponseWriter, errorMsg string) {
	w.Header().Set("Token-Error", errorMsg)
	w.WriteHeader(http.StatusUnauthorized)
}

func Challenge(w http.ResponseWriter) {
	w.Header().Set("www-authenticate", "bearer")
	w.WriteHeader(http.StatusUnauthorized)
}

type ResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.StatusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Hijack implements http.Hijacker interface for WebSocket support
func (rw *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, http.ErrNotSupported
	}
	return hijacker.Hijack()
}
