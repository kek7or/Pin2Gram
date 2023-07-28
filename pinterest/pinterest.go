package pinterest

type GetApiResponse struct {
	Status       string      `json:"status"`
	Code         int         `json:"code"`
	Message      string      `json:"message"`
	EndpointName string      `json:"endpoint_name"`
	Data         interface{} `json:"data"`
}
