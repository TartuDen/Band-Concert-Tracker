package visualization

import (
	"fmt"
	"groupie-tracker/get_info"
	"reflect"
	"testing"
)

func Test_ListOfLoc(t *testing.T) {

	relationData, _ := get_info.GetRelationData()
	result := ListOfLoc(relationData)
	expectedOutput := []string{
		"ANY", "argentina-buenos_aires", "argentina-la_plata", "argentina-san_isidro",
		"australia-brisbane", "australia-burswood", "australia-melbourne",
		"australia-new_south_wales", "australia-queensland", "australia-sydney",
		"australia-victoria", "australia-west_melbourne", "austria-klagenfurt",
		"austria-nickelsdorf", "austria-vienna", "belarus-minsk", "belgium-antwerp",
		"belgium-rotselaar", "belgium-werchter", "brazil-belo_horizonte",
		"brazil-brasilia", "brazil-porto_alegre", "brazil-recife",
		"brazil-rio_de_janeiro", "brazil-sao_paulo", "canada-montreal",
		"canada-quebec", "canada-toronto",
		"canada-vancouver", "canada-windsor", "chile-santiago", "china-changzhou",
		"china-hong_kong", "china-huizhou", "china-sanya", "colombia-bogota",
		"costa_rica-san_jose", "czechia-ostrava", "czechia-prague", "denmark-aalborg",
		"denmark-aarhus",
	}
	for i := 0; i < len(expectedOutput); i++ {
		if result[i] == expectedOutput[i] {
			fmt.Printf("Test Case %d: Passed\n", i+1)
			// fmt.Println(result[i])
		} else {
			t.Errorf("Test Case %d: Failed. Expected %s but got %s", i+1, expectedOutput[i], result[i])
		}
	}
}

func Test_CheckArtMembers(t *testing.T) {
	artData, _ := get_info.GetArtistData()
	testCases := []struct {
		description    string
		artistData     []get_info.Artist
		expectedMaxN   int
		expectedMaxOpt []int
	}{
		{
			description:    "Test of CheckArtMembers with example data",
			artistData:     artData,
			expectedMaxN:   8,
			expectedMaxOpt: []int{1, 2, 3, 4, 5, 6, 7, 8},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			resultMaxNExp, resultMaxOptExp := CheckArtMembers(testCase.artistData)

			if resultMaxNExp != testCase.expectedMaxN {
				t.Errorf("Test maxN: Failed. Expected %d but got %v", testCase.expectedMaxN, resultMaxNExp)
			}

			if !reflect.DeepEqual(resultMaxOptExp, testCase.expectedMaxOpt) {
				t.Errorf("Test maxOpt: Failed. Expected %v but got %v", testCase.expectedMaxOpt, resultMaxOptExp)
			}
		})
	}
}

func Test_CheckMinMaxCreationDate(t *testing.T) {
	artData, _ := get_info.GetArtistData()
	testCases := []struct {
		description             string
		artistData              []get_info.Artist
		expectedMinCreationDate int
		expectedMaxCreationDate int
	}{
		{description: "Test of CheckMinMaxCreationDate with example data",
			artistData:              artData,
			expectedMinCreationDate: 1958,
			expectedMaxCreationDate: 2015,
		}}
	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			resultMinCreationDate, resultMaxCreationDate := CheckMinMaxCreationDate(testCase.artistData)

			if resultMinCreationDate != testCase.expectedMinCreationDate {
				t.Errorf("Test maxN: Failed. Expected %d but got %v", testCase.expectedMinCreationDate, resultMinCreationDate)
			}

			if !reflect.DeepEqual(resultMaxCreationDate, testCase.expectedMaxCreationDate) {
				t.Errorf("Test maxOpt: Failed. Expected %v but got %v", testCase.expectedMaxCreationDate, resultMaxCreationDate)
			}
		})
	}
}

func Test_CheckMinMaxFirstAlbumDate(t *testing.T) {
	artData, _ := get_info.GetArtistData()
	testCases := []struct {
		description          string
		artistData           []get_info.Artist
		expectedMinAlbumDate int
		expectedMaxAlbumDate int
	}{
		{description: "Test of CheckMinMaxFirstAlbumDate with example data",
			artistData:           artData,
			expectedMinAlbumDate: 1963,
			expectedMaxAlbumDate: 2018,
		}}
	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			resultMinAlbumDate, resultMaxAlbumDate := CheckMinMaxFirstAlbumDate(testCase.artistData)

			if resultMinAlbumDate != testCase.expectedMinAlbumDate {
				t.Errorf("Test maxN: Failed. Expected %d but got %v", testCase.expectedMinAlbumDate, resultMinAlbumDate)
			}

			if !reflect.DeepEqual(resultMaxAlbumDate, testCase.expectedMaxAlbumDate) {
				t.Errorf("Test maxOpt: Failed. Expected %v but got %v", testCase.expectedMaxAlbumDate, resultMaxAlbumDate)
			}
		})
	}
}

func Test_ConvStringSliceToIntSlice(t *testing.T) {
	testCases := []struct {
		input    []string
		expected []int
	}{
		{[]string{"1", "2", "3"}, []int{1, 2, 3}},
		{[]string{"10", "20", "30"}, []int{10, 20, 30}},
		{[]string{"5", "15", "25"}, []int{5, 15, 25}},
		{[]string{"-1", "0", "1"}, []int{-1, 0, 1}},
	}

	for _, testCase := range testCases {
		result := ConvStringSliceToIntSlice(testCase.input)

		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Expected %v, but got %v for input %v", testCase.expected, result, testCase.input)
		}
	}
}

func Test_FilterFunc(t *testing.T) {

	artistsData, _ := get_info.GetArtistData()
	testCases := []struct {
		desc     string
		dataForF DataForFilterFunc

		artistsNameExpected []string
	}{
		{
			desc: "Test of FilterFunc with example data",
			dataForF: DataForFilterFunc{
				TeamNumFromValue: []int{0},
				MinAlbDate:       1963,
				MaxAlbDate:       2018,
				MaxCreationDate:  2000,
				MinCreationDate:  1995,
			},
			artistsNameExpected: []string{"SOJA", "Mamonas Assassinas", "Thirty Seconds to Mars", "Nickleback", "NWA", "Gorillaz", "Linkin Park", "Eminem", "Coldplay"},
		},
	}
	filteredResult, _, _ := FilterFunc(testCases[0].dataForF, artistsData)

	for idx, elem := range filteredResult {
		if elem.Name != testCases[0].artistsNameExpected[idx] {
			t.Errorf("T3.Expected %v, but got %v for input %v", testCases[0].artistsNameExpected[idx], elem.Name, testCases[0].artistsNameExpected)
		}
	}

}
