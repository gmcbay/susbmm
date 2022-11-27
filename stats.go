package main

type StatsResponse struct {
	Response struct {
		AllPvP struct {
			AllTime struct {
				KillsDeathsRatio struct {
					StatID string `json:"statId"`
					Basic  struct {
						Value        float64 `json:"value"`
						DisplayValue string  `json:"displayValue"`
					} `json:"basic"`
				} `json:"killsDeathsRatio"`
			} `json:"allTime"`
		} `json:"allPvP"`
	} `json:"Response"`
	ErrorCode       int    `json:"ErrorCode"`
	ThrottleSeconds int    `json:"ThrottleSeconds"`
	ErrorStatus     string `json:"ErrorStatus"`
	Message         string `json:"Message"`
	MessageData     struct {
	} `json:"MessageData"`
}
