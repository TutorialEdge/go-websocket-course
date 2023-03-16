package internal

// Event - "81' - Messi Scores a Volley"
type Event struct {
	// ChannelID - what channel the event corresponds to
	ChannelID string `json:"channel_id"`

	// What time the event occured
	Time string `json:"time"`
	// The associated message
	Message string `json:"message"`

	// we get the latest score
	HomeScore int `json:"home_score"`
	AwayScore int `json:"away_score"`
}
