package Controllers

import (
	m "UTS/Model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// GetAllRooms
func GetAllRooms(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT id, room_name FROM rooms"

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		// Handle error response
		sendErrorResponse(w, "Internal Server Error")
		return
	}
	var rooms []m.Room
	for rows.Next() {
		var room m.Room
		if err := rows.Scan(&room.ID, &room.RoomName); err != nil {
			log.Println(err)
			// Handle error response
			sendErrorResponse(w, "Internal Server Error")
			return
		}
		rooms = append(rooms, room)
	}
	w.Header().Set("Content-Type", "application/json")

	var response m.RoomsResponse
	response.Status = 200
	response.Message = "Success"
	response.Data.Rooms = rooms
	json.NewEncoder(w).Encode(response)
}

// GetRoomDetails
func GetRoomDetails(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	vars := mux.Vars(r)
	roomID := vars["room_id"]

	query := "SELECT id, room_name FROM rooms WHERE id = ?"
	row := db.QueryRow(query, roomID)

	var room m.Room
	if err := row.Scan(&room.ID, &room.RoomName); err != nil {
		log.Println(err)
		sendErrorResponse(w, "Room not found")
		return
	}

	queryParticipants := "SELECT p.id, p.id_account, a.username FROM participants p JOIN accounts a ON p.id_account = a.id WHERE p.id_room = ?"
	rows, err := db.Query(queryParticipants, roomID)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, "Internal Server Error")
		return
	}
	defer rows.Close()

	var participants []m.Participant
	for rows.Next() {
		var participant m.Participant
		if err := rows.Scan(&participant.ID, &participant.AccountID, &participant.Account.Username); err != nil {
			log.Println(err)
			sendErrorResponse(w, "Internal Server Error")
			return
		}
		participants = append(participants, participant)
	}

	w.Header().Set("Content-Type", "application/json")

	var response m.RoomDetailsResponse
	response.Status = 200
	response.Message = "Success"
	response.Data.Room = room
	response.Data.Participants = participants
	json.NewEncoder(w).Encode(response)
}

func InsertRoom(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendErrorResponse(w, "Bad Request")
		return
	}

	roomName := r.Form.Get("room_name")
	gameID := r.Form.Get("id_game")

	queryGame := "SELECT max_player FROM games WHERE id = ?"
	row := db.QueryRow(queryGame, gameID)

	var maxPlayer int
	if err := row.Scan(&maxPlayer); err != nil {
		log.Println(err)
		sendErrorResponse(w, "Internal Server Error")
		return
	}

	queryParticipantsCount := "SELECT COUNT(*) FROM participants WHERE id_room IN (SELECT id FROM rooms WHERE id_game = ?)"
	rowCount := db.QueryRow(queryParticipantsCount, gameID)

	var currentParticipantsCount int
	if err := rowCount.Scan(&currentParticipantsCount); err != nil {
		log.Println(err)
		sendErrorResponse(w, "Internal Server Error")
		return
	}

	if currentParticipantsCount >= maxPlayer {
		sendErrorResponse(w, "Room is full, cannot insert more participants")
		return
	}

	queryInsert := "INSERT INTO rooms (room_name, id_game) VALUES (?, ?)"
	result, err := db.Exec(queryInsert, roomName, gameID)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, "Internal Server Error")
		return
	}

	insertedRoomID, _ := result.LastInsertId()

	w.Header().Set("Content-Type", "application/json")

	var response m.SendResponse
	response.Status = 200
	response.Message = fmt.Sprintf("Room inserted successfully with ID %d", insertedRoomID)
	json.NewEncoder(w).Encode(response)
}

func sendSuccessResponse(w http.ResponseWriter, message string) {
	var response m.SendResponse
	response.Status = 200
	response.Message = message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func sendErrorResponse(w http.ResponseWriter, message string) {
	var response m.SendResponse
	response.Status = 400
	response.Message = message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
