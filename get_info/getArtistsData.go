package get_info

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Artist struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	LocDate      [][]string
	FinalRel     map[string][]string
}
type Relation struct {
	Id            int                 `json:"id"`
	DatesLocation map[string][]string `json:"datesLocations"`
}

type Relations struct {
	Index []Relation `json:"index"`
}

func GetArtistData() ([]Artist, error) {
	//Make an http get request to the api
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		fmt.Println("ERROR making get request:", err, "\nFrom: https://groupietrackers.herokuapp.com/api/artists")
		return nil, err
	}
	defer resp.Body.Close()

	//check the HTTP status code

	if resp.StatusCode != http.StatusOK {
		fmt.Println("HTTP request failed with status code: ", resp.StatusCode)
		return nil, err
	}

	//Decode the JSON response into a slice of Artists
	var artists []Artist
	err = json.NewDecoder(resp.Body).Decode(&artists)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, err
	}
	return artists, nil
}

func GetRelationData() ([]Relation, error) {
	var relations Relations
	//Make an http get request to the api
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		fmt.Println("ERROR making get request:", err, "\nFrom: https://groupietrackers.herokuapp.com/api/relation")
		return relations.Index, err
	}
	defer resp.Body.Close()

	//check the HTTP status code

	if resp.StatusCode != http.StatusOK {
		fmt.Println("HTTP request failed with status code: ", resp.StatusCode)
		return relations.Index, err
	}

	// Decode the JSON data into a slice of Relation structs
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&relations)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return relations.Index, err
	}

	return relations.Index, nil
}
