package visualization

import (
	"groupie-tracker/get_info"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

var SelectedConcertLocation = map[string][][]float64{}

var (
	maxCreationDate = 2023
	minCreationDate = 0
	minAlbumeDate   = 0
	maxAlbumeDate   = 2023
)

func HandlerMainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Initialize template
	templ, errTempl := template.New("Artists").Funcs(funcMap).ParseFS(viewTemplate, "htmlTemplates/*.html")
	if errTempl != nil {
		// Redirect to the error page
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	// Get data for the template
	viewData, err := getViewData(r)
	if err != nil {
		// Redirect to the error page
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	errTempl = templ.Execute(w, viewData)
	if errTempl != nil {
		// Redirect to the error page
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	}
}

// >>>>>>> Helper functions <<<<<<<

var funcMap = template.FuncMap{
	"add": func(a, b int) int {
		return a + b
	},
	"join": func(a []string) string {
		return strings.Join(a, ", ")
	},
	"sortDates": SortDates,
}

func getViewData(r *http.Request) (ViewData, []error) {
	//1. Get json data
	artistsData, errArt := get_info.GetArtistData()
	relationData, errRel := get_info.GetRelationData()

	combinedErr := []error{errArt, errRel}

	if errArt != nil || errRel != nil {
		return ViewData{}, combinedErr
	}

	//2. to get relevant data from API (json), to be used during pre-compilation of template
	maxMembersFromJson, maxMemberOptionFromJson := CheckArtMembers(artistsData)
	MinCreationDateFromJson, MaxCreationDateFromJson := CheckMinMaxCreationDate(artistsData)
	MinAlbDateFromJson, MaxAlbDateFromJson := CheckMinMaxFirstAlbumDate(artistsData)

	SlideMinDate := MinCreationDateFromJson
	SlideMaxDate := MaxCreationDateFromJson
	SlideMinAlbDate := MinAlbDateFromJson
	SlideMaxAlbDate := MaxAlbDateFromJson

	//3. to get filtered data from main page (in case of initial page, will return "0"s)
	minCreationYearFormValue, errMinY := strconv.Atoi(r.FormValue("minYear"))
	checkAtoiError(errMinY, "errMinY", minCreationYearFormValue)
	maxCreationYearFormValue, errMaxY := strconv.Atoi(r.FormValue("maxYear"))
	checkAtoiError(errMaxY, "errMaxY", maxCreationYearFormValue)
	teamNumFormValue := ConvStringSliceToIntSlice(r.Form["memberCount[]"])

	minYearAlbumFormValue, errMinYA := strconv.Atoi(r.FormValue("minYearAlbum"))
	checkAtoiError(errMinYA, "errMinYA", minYearAlbumFormValue)
	maxYearAlbumFormValue, errMaxYA := strconv.Atoi(r.FormValue("maxYearAlbum"))
	checkAtoiError(errMaxYA, "errMaxYA", maxYearAlbumFormValue)
	locFormValue := r.FormValue("location")

	searchQuery := r.FormValue("query")

	//to reset year filters
	if MinCreationDateFromJson != minCreationYearFormValue && minCreationYearFormValue != 0 {
		MinCreationDateFromJson = minCreationYearFormValue
	}
	if MaxCreationDateFromJson != maxCreationYearFormValue && maxCreationYearFormValue != 0 {
		MaxCreationDateFromJson = maxCreationYearFormValue
	}
	if MinAlbDateFromJson != minYearAlbumFormValue && minYearAlbumFormValue != 0 {
		MinAlbDateFromJson = minYearAlbumFormValue
	}
	if MaxAlbDateFromJson != maxYearAlbumFormValue && maxYearAlbumFormValue != 0 {
		MaxAlbDateFromJson = maxYearAlbumFormValue
	}

	//Filtering location
	for i, art := range artistsData {
		if art.Id == relationData[i].Id {
			sortedLoc := map[string][]string{}
			//for applied filter of locations
			if locFormValue != "ANY" && locFormValue != "" {
				for k, v := range relationData[i].DatesLocation {
					country := (strings.Split(k, "-"))[1] + "-" + (strings.Split(k, "-"))[0]
					if country == locFormValue {
						sortedLoc[k] = v
					}
				}
				artistsData[i].FinalRel = sortedLoc
				//for main page, printing all locations
			} else if locFormValue == "" {
				artistsData[i].FinalRel = relationData[i].DatesLocation
			} else if locFormValue == "ANY" {
				artistsData[i].FinalRel = relationData[i].DatesLocation
			}

		}
	}

	dataForF := DataForFilterFunc{
		TeamNumFromValue: teamNumFormValue,
		MinAlbDate:       MinAlbDateFromJson,
		MaxAlbDate:       MaxAlbDateFromJson,
		MaxCreationDate:  MaxCreationDateFromJson,
		MinCreationDate:  MinCreationDateFromJson,
		SearchQuery:      searchQuery,
	}
	//___________________________________________________
	filteredArtist, artistsCount, locCoords := FilterFunc(dataForF, artistsData)
	SelectedConcertLocation = locCoords
	//to compile final struct with filtered data
	viewData := ViewData{
		Artists:         filteredArtist,
		LocCoords:       locCoords,
		LenArtists:      artistsCount,
		Name:            "Groupie Tracker",
		MaxMembers:      maxMembersFromJson,
		MaxMemberOption: maxMemberOptionFromJson,
		MinCreationDate: SlideMinDate,
		MaxCreationDate: SlideMaxDate,
		MinYearAlbum:    SlideMinAlbDate,
		MaxYearAlbum:    SlideMaxAlbDate,
		ListOfLocs:      ListOfLoc(relationData),
		SelectedLoc:     "ANY",
	}
	return viewData, nil
}

func FilterFunc(data DataForFilterFunc, artistsData []get_info.Artist) ([]get_info.Artist, int, map[string][][]float64) {
	filteredArtist := []get_info.Artist{}
	LocCoords := map[string][][]float64{}
	//to save filtered data into a new struct "filteredArtist"
	var artistsCount int
	teamNumFormValue := data.TeamNumFromValue
	MinAlbDateFromJson := data.MinAlbDate
	MaxAlbDateFromJson := data.MaxAlbDate
	MaxCreationDateFromJson := data.MaxCreationDate
	MinCreationDateFromJson := data.MinCreationDate
	SearchQuery := strings.ToLower(data.SearchQuery)
	checkInitialPage := true
	for _, artistInfo := range artistsData {
		NameToSearch := strings.ToLower(artistInfo.Name)
		curFAD, errCurFAD := strconv.Atoi((strings.Split(artistInfo.FirstAlbum, "-"))[2])
		checkAtoiError(errCurFAD, "errCurFAD", curFAD)
		if len(teamNumFormValue) == 0 && len(SearchQuery) == 0 { //creating initial page
			filteredArtist = artistsData
			artistsCount = len(artistsData)
			break
		} else if len(teamNumFormValue) >= 1 { //creating filtered pages
			for _, nubMemb := range teamNumFormValue {
				if curFAD >= MinAlbDateFromJson && curFAD <= MaxAlbDateFromJson && artistInfo.CreationDate >= MinCreationDateFromJson &&
					artistInfo.CreationDate <= MaxCreationDateFromJson &&
					(len(artistInfo.Members) == nubMemb || nubMemb == 0) &&
					(len(artistInfo.FinalRel) != 0) &&
					len(SearchQuery) == 0 {
					filteredArtist = append(filteredArtist, artistInfo)
					artistsCount++
					checkInitialPage = false
				}

			}

		} else if (SearchQuery == NameToSearch) || (strings.Contains(NameToSearch, SearchQuery)) {
			filteredArtist = append(filteredArtist, artistInfo)
			artistsCount++
			checkInitialPage = false

		}
	}
	if !checkInitialPage {
		LocCoords = GeoData(filteredArtist)
	} else {
		LocCoords = nil
	}
	return filteredArtist, artistsCount, LocCoords
}

func ConvStringSliceToIntSlice(data []string) []int {
	newSlice := []int{}
	for _, n := range data {
		tempNumb, _ := strconv.Atoi(n)
		newSlice = append(newSlice, tempNumb)
	}
	return newSlice
}

func ListOfLoc(data []get_info.Relation) []string {
	// map[string]struct{} to keep track of unique locations
	uniqueLocations := make(map[string]struct{})

	for _, band := range data {
		for loc := range band.DatesLocation {
			parts := strings.Split(loc, "-")
			if len(parts) == 2 {
				// Swap the order of parts and join them back
				joinedLoc := parts[1] + "-" + parts[0]
				uniqueLocations[joinedLoc] = struct{}{}
			}
		}
	}
	uniqueLocations["ANY"] = struct{}{}

	sortedLocations := make([]string, 0, len(uniqueLocations))
	for loc := range uniqueLocations {
		sortedLocations = append(sortedLocations, loc)
	}

	sort.Strings(sortedLocations)
	return sortedLocations
}

func CheckArtMembers(data []get_info.Artist) (int, []int) {
	maxN := 0
	for _, band := range data {
		if len(band.Members) > maxN {
			maxN = len(band.Members)
		}

	}
	maxOpt := []int{}
	for i := 1; i <= maxN; i++ {
		maxOpt = append(maxOpt, i)
	}
	return maxN, maxOpt
}

func CheckMinMaxCreationDate(data []get_info.Artist) (int, int) {
	minCD := maxCreationDate
	maxCD := minCreationDate
	for _, band := range data {
		if band.CreationDate <= minCD {
			minCD = band.CreationDate
		}
		if band.CreationDate >= maxCD {
			maxCD = band.CreationDate
		}
	}
	return minCD, maxCD
}

func CheckMinMaxFirstAlbumDate(data []get_info.Artist) (int, int) {
	minAlbDat := maxAlbumeDate
	maxAlbDat := minAlbumeDate
	for _, band := range data {
		tempToStoreYear, errTempToStoreY := strconv.Atoi((strings.Split(band.FirstAlbum, "-"))[2])
		checkAtoiError(errTempToStoreY, "errTempToStoreY", tempToStoreYear)

		if tempToStoreYear <= minAlbDat {
			minAlbDat = tempToStoreYear
		}
		if tempToStoreYear >= maxAlbDat {
			maxAlbDat = tempToStoreYear
		}
	}
	return minAlbDat, maxAlbDat
}

func checkAtoiError(err error, errName string, varInt int) {
	if varInt != 0 && err != nil {
		log.Fatal(errName, err)
	}
}
