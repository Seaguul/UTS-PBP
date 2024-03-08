package models

type Account struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type Game struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	MaxPlayers int    `json:"max_players"`
}

type Room struct {
	ID       int    `json:"id"`
	RoomName string `json:"room_name"`
	GameID   int    `json:"game_id"`
}

type Participant struct {
	ID        int     `json:"id"`
	RoomID    int     `json:"room_id"`
	AccountID int     `json:"account_id"`
	Account   Account `json:"account"` // Include Account information
}

type SendResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type AccountsResponse struct {
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Data    []Account `json:"data"`
}

type GamesResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []Game `json:"data"`
}

type ParticipantsResponse struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Data    []Participant `json:"data"`
}

type RoomsResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Rooms []Room `json:"rooms"`
	} `json:"data"`
}

type RoomDetailsResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Room         Room          `json:"room"`
		Participants []Participant `json:"participants"`
	} `json:"data"`
}
