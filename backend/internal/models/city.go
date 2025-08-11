package models

type CityConfig struct {
	Name   string  `json:"name"`
	Bounds *Bounds `json:"Bounds"`
}

type Bounds struct {
	Northeast Coordinate `json:"northeast"`
	Southwest Coordinate `json:"southwest"`
}

type Coordinate struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
