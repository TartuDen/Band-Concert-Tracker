package visualization

import (
	"net/http"
	"os"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	// Read and serve the error page HTML
	errorPageHTML, err := os.ReadFile("visualization/errorPage/error.html")
	if err != nil {
		http.Error(w, "Error page not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(errorPageHTML)
}
