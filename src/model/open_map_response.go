package model

type OpenMapResponse struct {
	Main Main `json:"main"`
	Wind Wind `json:"wind"`
}

type Main struct {
	Temperature float64 `json:"temp"`
}

type Wind struct {
	WindSpeed float64 `json:"speed"`
}
