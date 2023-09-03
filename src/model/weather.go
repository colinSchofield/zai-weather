package model

type Weather struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    *Data  `json:"data,omitempty"`
}

type Data struct {
	Temperature int `json:"temperature_degrees"`
	WindSpeed   int `json:"wind_speed"`
}
