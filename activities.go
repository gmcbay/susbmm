package main

import (
	"time"
)

type ActivitiesData struct {
	Activities []struct {
		Period          time.Time `json:"period"`
		ActivityDetails struct {
			InstanceID string `json:"instanceId"`
		} `json:"activityDetails"`
	} `json:"activities"`
}

type ActivitiesDataResponse struct {
	Response        ActivitiesData `json:"Response"`
	ErrorCode       int            `json:"ErrorCode"`
	ThrottleSeconds int            `json:"ThrottleSeconds"`
	ErrorStatus     string         `json:"ErrorStatus"`
	Message         string         `json:"Message"`
	MessageData     struct {
	} `json:"MessageData"`
}
