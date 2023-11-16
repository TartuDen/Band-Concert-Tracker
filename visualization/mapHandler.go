package visualization

import (
	"fmt"
	"html/template"
	"net/http"
)

func MapHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/map" {
		http.NotFound(w, r)
		return
	}

	// Initialize template
	templ, errTempl := template.New("Maps").ParseFS(viewMaps, "mapTempl/*.html")
	if errTempl != nil {
		fmt.Println("errTempl - in formting templ", errTempl)
		// Redirect to the error page
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}
	initialLoc := [][]float64{}
	for _, sliceCoords := range SelectedConcertLocation {
		initialLoc = append(initialLoc, sliceCoords[0])
		break
	}
	// Get data for the template
	MapViewData := MapStruct{
		ConcertData: SelectedConcertLocation,
		InitialLoc:  initialLoc,
	}

	errTempl = templ.Execute(w, MapViewData)
	if errTempl != nil {
		fmt.Println("errTempl", errTempl)
		// Redirect to the error page
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	}
}
