package pinterest

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type GetApiResponse struct {
	Status       string `json:"status"`
	Code         int    `json:"code"`
	Message      string `json:"message"`
	EndpointName string `json:"endpoint_name"`
	Data         Data   `json:"data"`
}

type Data struct {
	User  User  `json:"user"`
	Pins  []Pin `json:"pins"`
	Board Board `json:"board"`
}

type User struct {
	FollowerCount int    `json:"follower_count"`
	ID            string `json:"id"`
	ImageSmallURL string `json:"image_small_url"`
	FullName      string `json:"full_name"`
	About         string `json:"about"`
	ProfileURL    string `json:"profile_url"`
	PinCount      int    `json:"pin_count"`
}

type Pin struct {
	Attribution       interface{}       `json:"attribution"`
	DominantColor     string            `json:"dominant_color"`
	StoryPinData      interface{}       `json:"story_pin_data"`
	Description       string            `json:"description"`
	NativeCreator     User              `json:"native_creator"`
	ID                string            `json:"id"`
	Images            Images            `json:"images"`
	RepinCount        int               `json:"repin_count"`
	Domain            string            `json:"domain"`
	AggregatedPinData AggregatedPinData `json:"aggregated_pin_data"`
	Link              interface{}       `json:"link"`
	Pinner            User              `json:"pinner"`
	IsVideo           bool              `json:"is_video"`
	Embed             interface{}       `json:"embed"`
}

type Images struct {
	I237x Image `json:"237x"`
	I564x Image `json:"564x"`
}

type Image struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
}

type AggregatedPinData struct {
	AggregatedStats AggregatedStats `json:"aggregated_stats"`
}

type AggregatedStats struct {
	Saves int `json:"saves"`
	Done  int `json:"done"`
}

type Board struct {
	FollowerCount     int    `json:"follower_count"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	ImageThumbnailURL string `json:"image_thumbnail_url"`
	URL               string `json:"url"`
	ID                string `json:"id"`
	PinCount          int    `json:"pin_count"`
}

func GetPinsFromBoard(board string) ([]Pin, error) {
	uri := fmt.Sprintf(PIDGETS_URI, board)
	resp, err := http.Get(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to access '%s' board", board)
	}
	defer resp.Body.Close()

	var decoded GetApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return nil, fmt.Errorf("failed to decode '%s' board info to json", board)
	}

	if decoded.Status != "success" {
		return nil, fmt.Errorf("failed to get info about '%s' board", board)
	}

	return decoded.Data.Pins, nil
}
