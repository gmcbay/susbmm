package main

type ClassType int

const (
	Titan   ClassType = 0
	Hunter  ClassType = 1
	Warlock ClassType = 2
)

func (ct ClassType) String() string {
	switch ct {
	case Titan:
		return "Titan"

	case Hunter:
		return "Hunter"

	case Warlock:
		return "Warlock"

	default:
		return "Unknown"
	}
}

type CharacterData struct {
	MembershipID   string    `json:"membershipId"`
	MembershipType int       `json:"membershipType"`
	CharacterID    string    `json:"characterId"`
	ClassType      ClassType `json:"classType"`
}

func getCharacterClassName(characterData CharacterData) string {
	return characterData.ClassType.String()
}
