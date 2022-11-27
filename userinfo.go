package main

type UserInfoCard struct {
	MembershipType              int    `json:"membershipType"`
	MembershipID                string `json:"membershipId"`
	BungieGlobalDisplayName     string `json:"bungieGlobalDisplayName"`
	BungieGlobalDisplayNameCode int    `json:"bungieGlobalDisplayNameCode"`
}

type UserInfoCardArrayResponse struct {
	Response        []UserInfoCard `json:"Response"`
	ErrorCode       int            `json:"ErrorCode"`
	ThrottleSeconds int            `json:"ThrottleSeconds"`
	ErrorStatus     string         `json:"ErrorStatus"`
	Message         string         `json:"Message"`
	MessageData     struct {
	} `json:"MessageData"`
}
