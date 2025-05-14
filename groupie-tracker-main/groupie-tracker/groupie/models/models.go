package models

import (
	"encoding/json"
	"net/http"
)

func ConvertToStringSlice(v interface{}) []string {
	switch val := v.(type) {
	case string:
		return []string{val}
	case []interface{}:
		var result []string
		for _, item := range val {
			if str, ok := item.(string); ok {
				result = append(result, str)
			}
		}
		return result
	default:
		return nil
	}
}

type LocationsData struct {
	Index []struct {
		ID        int      `json:"id"`
		Locations []string `json:"locations"`
	} `json:"index"`
}

func FetchLocations() (map[int]map[string][]string, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rawJSON map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&rawJSON)
	if err != nil {
		return nil, err
	}

	locationMap := make(map[int]map[string][]string)
	if indexArray, ok := rawJSON["index"].([]interface{}); ok {
		for _, item := range indexArray {
			if locItem, ok := item.(map[string]interface{}); ok {
				id := int(locItem["id"].(float64))
				locationsMap := make(map[string][]string)

				if locations, ok := locItem["locations"].([]interface{}); ok {
					cities := make([]string, 0, len(locations))
					for _, loc := range locations {
						if cityStr, ok := loc.(string); ok {
							cities = append(cities, cityStr)
						}
					}
					locationsMap["cities"] = cities
				}

				locationMap[id] = locationsMap
			}
		}
	}

	return locationMap, nil
}

type Artist struct {
	ID           int                 `json:"id"`
	Image        string              `json:"image"`
	Name         string              `json:"name"`
	Members      []string            `json:"members"`
	CreationDate int                 `json:"creationDate"`
	FirstAlbum   string              `json:"firstAlbum"`
	Locations    map[string][]string `json:"locations"`
	ConcertDates map[string][]string `json:"concertDates"`
}

func FetchArtists() ([]Artist, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rawArtists []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&rawArtists)
	if err != nil {
		return nil, err
	}

	locations, err := FetchLocations()
	if err != nil {
		locations = make(map[int]map[string][]string)
	}

	concertDates, err := FetchConcertDates()
	if err != nil {
		concertDates = make(map[int]map[string][]string)
	}

	var artists []Artist
	for _, raw := range rawArtists {
		id := int(raw["id"].(float64))
		artistLocations := locations[id]
		if artistLocations == nil {
			artistLocations = make(map[string][]string)
		}

		artistDates := concertDates[id]
		if artistDates == nil {
			artistDates = make(map[string][]string)
		}

		artist := Artist{
			ID:           id,
			Image:        raw["image"].(string),
			Name:         raw["name"].(string),
			Members:      ConvertToStringSlice(raw["members"]),
			CreationDate: int(raw["creationDate"].(float64)),
			FirstAlbum:   raw["firstAlbum"].(string),
			Locations:    artistLocations,
			ConcertDates: artistDates,
		}
		artists = append(artists, artist)
	}
	return artists, nil
}

type RelationData struct {
	Index []struct {
		ID             int                 `json:"id"`
		DatesLocations map[string][]string `json:"datesLocations"`
	} `json:"index"`
}

func FetchConcertDates() (map[int]map[string][]string, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rawJSON map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&rawJSON)
	if err != nil {
		return nil, err
	}

	concertMap := make(map[int]map[string][]string)

	if indexArray, ok := rawJSON["index"].([]interface{}); ok {
		for _, item := range indexArray {
			if concertItem, ok := item.(map[string]interface{}); ok {
				id := int(concertItem["id"].(float64))
				datesMap := make(map[string][]string)

				if datesLoc, ok := concertItem["datesLocations"].(map[string]interface{}); ok {
					for location, dates := range datesLoc {
						if datesArray, ok := dates.([]interface{}); ok {
							dateStrings := make([]string, 0, len(datesArray))
							for _, date := range datesArray {
								if dateStr, ok := date.(string); ok {
									dateStrings = append(dateStrings, dateStr)
								}
							}
							datesMap[location] = dateStrings
						}
					}
				}

				concertMap[id] = datesMap
			}
		}
	}

	return concertMap, nil
}

func FetchArtistByID(id int) (*Artist, error) {
	artists, err := FetchArtists()
	if err != nil {
		return nil, err
	}

	for _, artist := range artists {
		if artist.ID == id {
			return &artist, nil
		}
	}

	return nil, nil
}
