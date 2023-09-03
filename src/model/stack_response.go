package model

type StackResponse struct {
	Error   Error   `json:"error"`
	Current Current `json:"current"`
}

type Error struct {
	Code int `json:"code"`
}

type Current struct {
	Temperature int `json:"temperature"`
	WindSpeed   int `json:"wind_speed"`
}
