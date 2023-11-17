package visualization

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/get_info"
	"io"
	"net/http"
	"strings"
)

func GeoData(artists []get_info.Artist) map[string][][]float64 {
	locations := map[string][]string{}
	for _, data := range artists {
		tempSlice := []string{}
		for loc, _ := range data.FinalRel {
			tempSlice = append(tempSlice, loc)
		}
		locations[data.Name] = tempSlice

	}

	locationsCoords := map[string][][]float64{}
	for bandName, sliceLocs := range locations {
		tempSliceCoords := make([][]float64, len(sliceLocs)) // Initialize the inner slices

		for idxLoc, loc := range sliceLocs {
			city := strings.Split(loc, "-")[0]
			country := strings.Split(loc, "-")[1]
			lat, lng := reverseGeocode(fmt.Sprintf("%s,%s", city, country))

			// Append to the inner slice
			tempSliceCoords[idxLoc] = append(tempSliceCoords[idxLoc], lat)
			tempSliceCoords[idxLoc] = append(tempSliceCoords[idxLoc], lng)
		}
		locationsCoords[bandName] = tempSliceCoords
	}
	return locationsCoords
}

type GeocodingResponseAdress struct {
	Results []struct {
		Geometry struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
	Status string `json:"status"`
}

// getting actual coordinates from addresses using Map API
func reverseGeocode(address string) (float64, float64) {
	url := "https://maps.googleapis.com/maps/api/geocode/json?address=" + address + "&key=" + apiKey

	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error:", err)
		// return
	}
	var geocodingResponse GeocodingResponseAdress

	err = json.Unmarshal(body, &geocodingResponse)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		// return
	}
	var lat, lng float64
	// fmt.Println("Status:", geocodingResponse.Status)
	if len(geocodingResponse.Results) > 0 {
		lat = geocodingResponse.Results[0].Geometry.Location.Lat
		lng = geocodingResponse.Results[0].Geometry.Location.Lng
		// fmt.Printf("Latitude: %f, Longitude: %f\n", lat, lng)
		// fmt.Printf("https://www.google.com/maps/search/?api=1&query=%f,%f\n", lat, lng)
	}
	return lat, lng
}
