package routes

import (
	"html/template"
	"net/http"

	"github.com/sergivb01/udrop-leaderboards/leaderboard"
)

// Templates defines a list of templates
var Templates *template.Template

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if err := Templates.ExecuteTemplate(w, "index", map[string]interface{}{
		"Players": genRandomPlayers(),
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func genRandomPlayers() []leaderboard.Player {
	var players []leaderboard.Player

	for i := 0; i < 20; i++ {
		players = append(players, leaderboard.Player{
			Playername: "Steve",
			UUID:       "hello",
			Stats:      nil,
		})
	}

	return players
}
