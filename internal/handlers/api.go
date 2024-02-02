package handlers

import (
	"fmt"
	"net/http"
)

func ServeHelloWorldAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<div> Hello World!</div>")
}
