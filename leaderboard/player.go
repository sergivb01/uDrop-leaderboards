package leaderboard

type Player struct {
	Playername string                 `json:"playername"`
	UUID       string                 `json:"uuid"`
	Stats      map[string]interface{} `json:"statistics"`
}
