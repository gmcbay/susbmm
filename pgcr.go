package main

type PGCRData struct {
	Entries []struct {
		Player struct {
			DestinyUserInfo struct {
				MembershipType              int    `json:"membershipType"`
				MembershipID                string `json:"membershipId"`
				BungieGlobalDisplayName     string `json:"bungieGlobalDisplayName"`
				BungieGlobalDisplayNameCode int    `json:"bungieGlobalDisplayNameCode"`
			} `json:"destinyUserInfo"`
		} `json:"player"`
		Values struct {
			Kills struct {
				Basic struct {
					Value        float64 `json:"value"`
					DisplayValue string  `json:"displayValue"`
				} `json:"basic"`
			} `json:"kills"`
			Team struct {
				Basic struct {
					Value        float64 `json:"value"`
					DisplayValue string  `json:"displayValue"`
				} `json:"basic"`
			} `json:"team"`
		} `json:"values"`
	} `json:"entries"`
}

type PGCRDataResponse struct {
	Response        PGCRData `json:"Response"`
	ErrorCode       int      `json:"ErrorCode"`
	ThrottleSeconds int      `json:"ThrottleSeconds"`
	ErrorStatus     string   `json:"ErrorStatus"`
	Message         string   `json:"Message"`
	MessageData     struct {
	} `json:"MessageData"`
}
