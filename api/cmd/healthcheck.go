package main

import (
	"fmt"
	"net/http"
)

func (app *application) Healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "status: healthy")
	fmt.Fprintf(w, "env: %s\n", app.config.env)
	fmt.Fprintf(w, "version: %s\n", version)
}
