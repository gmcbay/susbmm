package main

type ProfileData struct {
	Profile struct {
		Data struct {
			UserInfo     UserInfoCard `json:"userInfo"`
			CharacterIds []string     `json:"characterIds"`
		} `json:"data"`
		Privacy int `json:"privacy"`
	} `json:"profile"`
	Characters struct {
		CharacterData map[string]CharacterData `json:"data"`
	} `json:"characters"`
}

type ProfileDataResponse struct {
	Response        ProfileData `json:"Response"`
	ErrorCode       int         `json:"ErrorCode"`
	ThrottleSeconds int         `json:"ThrottleSeconds"`
	ErrorStatus     string      `json:"ErrorStatus"`
	Message         string      `json:"Message"`
	MessageData     struct {
	} `json:"MessageData"`
}
