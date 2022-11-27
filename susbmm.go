package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

func main() {
	fmt.Printf("API_KEY: %+v\n", API_KEY)
	
	args := strings.Join(os.Args[1:], " ")

	splitName := strings.Split(args, "#")

	if len(splitName) != 2 {
		fmt.Printf("Usage: %s [BungieName]#[BungieCode]\n", os.Args[0])
		os.Exit(-1)
	}

	fmt.Printf("Looking for: %s#%s\n", splitName[0], splitName[1])

	userInfo := getDestinyPlayerByName(splitName[0], splitName[1])

	if userInfo.ErrorCode != ErrSuccess {
		if userInfo.ErrorCode == ErrApiKeyInvalid {
			log.Fatalf("You must request a personal API key from bungie and replace the value for API_KEY in api.go before building this tool\n")
		}

		log.Fatalf(
			"Unexpected error response while getting Destiny Player Name: %d\n",
			userInfo.ErrorCode)
	}

	if len(userInfo.Response) < 1 {
		log.Fatalf("User not found: %s#%s\n", splitName[0], splitName[1])
	}

	fmt.Printf("Found player: %v#%v\n",
		userInfo.Response[0].BungieGlobalDisplayName,
		userInfo.Response[0].BungieGlobalDisplayNameCode)

	charProfileData := getCharacterInfo(userInfo)

	instanceIDs := make([]string, 0, 2048)

	for _, character := range charProfileData.Response.Characters.CharacterData {
		characterClassName := getCharacterClassName(character)
		fmt.Printf("Found character: %s\n", characterClassName)
		fmt.Printf("Fetching Control/SBMM activity list for %s\n",
			characterClassName)
		instanceIDs = getControlSbmmActivitiesForCharacter(character,
			instanceIDs)
	}

	pgcrCount := len(instanceIDs)

	fmt.Printf("Found %d Control/SBMM games\n", pgcrCount)

	var wg sync.WaitGroup

	smurfChan := make(chan string, pgcrCount*24)

	count := 0
	startBatchTime := int64(time.Now().UnixMilli())

	for _, instanceID := range instanceIDs {
		wg.Add(1)

		go getPossibleSmurfsFromPGCR(smurfChan, &wg,
			userInfo.Response[0].MembershipID, instanceID)

		count++

		if count%50 == 0 {
			fmt.Printf("PGCR requests: %+v of %+v\n", count, pgcrCount)

			wg.Wait()

			endBatchTime := int64(time.Now().UnixMilli())

			delta := endBatchTime - startBatchTime
			waitTime := 2000 - delta

			if delta < 2000 {
				time.Sleep(time.Duration(waitTime) * time.Millisecond)
			}

			startBatchTime = int64(time.Now().UnixMilli())
		}
	}

	wg.Wait()
	close(smurfChan)

	smurfCounts := make(map[string]int)

	for smurfId := range smurfChan {
		count, ok := smurfCounts[smurfId]

		if ok {
			count++
		} else {
			count = 0
		}

		smurfCounts[smurfId] = count
	}

	keys := make([]string, 0, len(smurfCounts))

	for key := range smurfCounts {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return smurfCounts[keys[i]] > smurfCounts[keys[j]]
	})

	for i := 0; i < len(keys); i++ {
		count := smurfCounts[keys[i]]

		if count >= SMURF_GAME_APPEARANCE_FLOOR {
			smurfSplitName := strings.Split(keys[i], "#")
			smurfInfo := getDestinyPlayerByName(smurfSplitName[0],
				smurfSplitName[1])
			smurfKD := getKD(smurfInfo)

			if smurfKD <= SMURF_KD_FLOOR {
				fmt.Printf(
					"Possible smurf associated with this player: %+v (%+v KD)\n",
					keys[i], smurfKD)
			}
		}
	}

	fmt.Printf("Done...\n")
}
