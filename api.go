package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const API_KEY = "PASTE_YOUR_API_KEY"

// The most amount of kills a smurf can get in a game,
// if they score more than this that game isn't counted
// when trying to detect smurfs
const SMURF_SINGLE_GAME_KILL_FLOOR = 2
const SMURF_GAME_APPEARANCE_FLOOR = 3
const SMURF_KD_FLOOR = 0.3

const API_ROOT = "https://www.bungie.net/Platform"
const ACTIVITIES_COUNT_MAX = 250
const QP_CONTROL_MODE = 19

const (
	ErrNone          int = 0
	ErrSuccess       int = 1
	ErrApiKeyInvalid int = 2101
)

var client = &http.Client{}
var sbmmStart time.Time

// define sbmmStart as season of the deep start.
// games before that cut-off won't be looked at
func init() {
	var err error

	sbmmStart, err = time.Parse(time.RFC3339, "2023-05-23T17:00:00Z")

	if err != nil {
		log.Fatal(err)
	}
}

func makeApiRequest(method string, url string, body string) []byte {
	request, err := http.NewRequest(method, url, strings.NewReader(body))

	if err != nil {
		log.Fatal(err)
	}

	request.Header.Add("X-API-Key", API_KEY)

	response, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
	}

	responseBody, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	return responseBody
}

func getDestinyPlayerByName(displayName string,
	displayNameCode string) (userInfos UserInfoCardArrayResponse) {
	name := &PlayerName{displayName, displayNameCode}

	postData, err := json.Marshal(name)

	body := makeApiRequest("POST", fmt.Sprintf(
		"%+v/Destiny2/SearchDestinyPlayerByBungieName/-1/",
		API_ROOT), string(postData))

	err = json.Unmarshal(body, &userInfos)

	if err != nil {
		log.Fatal(err)
	}

	if userInfos.ErrorCode != ErrSuccess {
		log.Fatalf("%+v: %+v\n", userInfos.ErrorStatus, userInfos.Message)
	}

	return
}

func getKD(userInfos UserInfoCardArrayResponse) float64 {
	body := makeApiRequest("GET", fmt.Sprintf(
		"%+v/Destiny2/%d/Account/%s/Character/0/Stats/",
		API_ROOT, userInfos.Response[0].MembershipType,
		userInfos.Response[0].MembershipID), "")

	statData := StatsResponse{}
	err := json.Unmarshal(body, &statData)

	if err != nil {
		log.Fatal(err)
	}

	if statData.ErrorCode != ErrSuccess {
		log.Fatalf("%+v: %+v\n", statData.ErrorStatus, statData.Message)
	}

	return statData.Response.AllPvP.AllTime.KillsDeathsRatio.Basic.Value
}

func getCharacterInfo(
	userInfos UserInfoCardArrayResponse) (profileData ProfileDataResponse) {
	body := makeApiRequest("GET", fmt.Sprintf(
		"%+v/Destiny2/%d/Profile/%s/?components=100,200",
		API_ROOT, userInfos.Response[0].MembershipType,
		userInfos.Response[0].MembershipID), "")

	err := json.Unmarshal(body, &profileData)

	if err != nil {
		log.Fatal(err)
	}

	if profileData.ErrorCode != ErrSuccess {
		log.Fatalf("%+v: %+v\n", profileData.ErrorStatus, profileData.Message)
	}

	return
}

func getControlSbmmActivitiesForCharacter(characterData CharacterData,
	instanceIDs []string) []string {
	return getControlSbmmActivitiesPageForCharacter(characterData, 0,
		instanceIDs)
}

func getControlSbmmActivitiesPageForCharacter(characterData CharacterData,
	page int, instanceIDs []string) []string {
	reqStartTime := int64(time.Now().UnixMilli())

	body := makeApiRequest("GET", fmt.Sprintf(
		"%+v/Destiny2/%d/Account/%s/Character/%s/Stats/Activities/?mode=%d&count=%d&page=%d",
		API_ROOT, characterData.MembershipType, characterData.MembershipID,
		characterData.CharacterID, QP_CONTROL_MODE, ACTIVITIES_COUNT_MAX, page),
		"")

	activitiesData := ActivitiesDataResponse{}
	err := json.Unmarshal(body, &activitiesData)

	if err != nil {
		log.Fatal(err)
	}

	if activitiesData.ErrorCode != ErrSuccess {
		log.Fatalf("%+v: %+v\n", activitiesData.ErrorStatus,
			activitiesData.Message)
	}

	activityCount := len(activitiesData.Response.Activities)

	reqEndTime := int64(time.Now().UnixMilli())

	reqTime := (reqEndTime - reqStartTime)

	if reqTime < 50 {
		time.Sleep(time.Duration(50-reqTime) * time.Millisecond)
	}

	for i := 0; i < activityCount; i++ {
		// when we start seeing dates prior to sbmmStart, we're done for
		// this set of activities, return what we've found
		if sbmmStart.After(activitiesData.Response.Activities[i].Period) {
			return instanceIDs
		}

		instanceIDs = append(instanceIDs,
			activitiesData.Response.Activities[i].ActivityDetails.InstanceID)
	}

	// if we pulled less activities than we requested, we're done for
	// this set of activities, return what we've found
	if activityCount < ACTIVITIES_COUNT_MAX {
		return instanceIDs
	}

	// if not done, go to the next page
	return getControlSbmmActivitiesPageForCharacter(characterData, page+1,
		instanceIDs)
}

func getPossibleSmurfsFromPGCR(ch chan string, wg *sync.WaitGroup,
	membershipID string, instanceID string) {
	defer wg.Done()

	body := makeApiRequest("GET", fmt.Sprintf(
		"%+v/Destiny2/Stats/PostGameCarnageReport/%v/",
		API_ROOT, instanceID), "")

	pgcrData := PGCRDataResponse{}
	err := json.Unmarshal(body, &pgcrData)

	if err != nil {
		log.Fatal(err)
	}

	if pgcrData.ErrorCode != ErrSuccess {
		log.Fatalf("%+v: %+v\n", pgcrData.ErrorStatus, pgcrData.Message)
	}

	teamID := -1

	for _, entry := range pgcrData.Response.Entries {
		if membershipID == entry.Player.DestinyUserInfo.MembershipID {
			teamID = int(entry.Values.Team.Basic.Value)
			break
		}
	}

	if teamID == -1 {
		teamID = 0
	}

	for _, entry := range pgcrData.Response.Entries {
		compareTeamID := int(entry.Values.Team.Basic.Value)

		if (membershipID != entry.Player.DestinyUserInfo.MembershipID) &&
			(teamID == compareTeamID || compareTeamID == -1) {
			killCount := int(entry.Values.Kills.Basic.Value)

			if killCount <= SMURF_SINGLE_GAME_KILL_FLOOR &&
				len(entry.Player.DestinyUserInfo.BungieGlobalDisplayName) > 0 {
				ch <- fmt.Sprintf("%s#%04d",
					entry.Player.DestinyUserInfo.BungieGlobalDisplayName,
					entry.Player.DestinyUserInfo.BungieGlobalDisplayNameCode)
			}
		}
	}
}
