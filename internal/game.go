package internal

// Channels - represents all the various channels
// you can click into.
type Channels struct {
	// a map of channel-ids to games
	Games map[string]Game
}

// Game - represents a singular game
type Game struct {
	// Celtic
	HomeTeam string
	// Pollok FC
	AwayTeam string

	// 0-0, 1-0, 2-3
	HomeTeamScore int
	AwayTeamScore int

	// 0-90 minutes
	GameTime string
}

// Event - "81' - Messi Scores a Volley"
type Event struct {
	// What time the event occured
	Time string
	// The associated message
	Message string

	// we get the latest score
	HomeScore int
	AwayScore int
}
