package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/notnil/chess"
)

type Game struct {
	lastPlayedMoveID *int   // optional
	gameState        string // PGN or something
}

type DBType *map[int]*Game

/*
In memory 'database'
*/
func registerClient(botID int, db DBType) error {
	if _, ok := (*db)[botID]; ok {
		// contains botID already
		return fmt.Errorf("contains botID already")
	} else {
		(*db)[botID] = nil
		return nil
	}
}

func startGame(botID int, db DBType) error {
	if _, ok := (*db)[botID]; !ok {
		// doesn't contain botID
		return fmt.Errorf("bot id doesn't exist")
	}
	// check if a another game is running
	if game := (*db)[botID]; game != nil {
		return fmt.Errorf("another game is running already")
	}
	(*db)[botID] = &Game{nil, ""}
	return nil
}

func playMove(botID int, playerID int, move string, db DBType) error {
	if _, ok := (*db)[botID]; !ok {
		// doesn't contain botID
		return fmt.Errorf("bot id doesn't exist")
	}
	// check if a game is running
	if game := (*db)[botID]; game == nil {
		return fmt.Errorf("no game running")
	}

	// check if we are not making the a new move with the same player
	// that played the previous move
	if game := (*db)[botID]; game.lastPlayedMoveID != nil && *game.lastPlayedMoveID == playerID {
		return fmt.Errorf("consecutive moves detected")
	}

	// instantiate a board and load current state
	pgn, err := chess.PGN(strings.NewReader((*db)[botID].gameState))
	if err != nil {
		return err
	}
	game := chess.NewGame(pgn)

	// try making the move
	err = game.MoveStr(move)
	if err != nil {
		return err
	}

	// update the board & player
	(*db)[botID].gameState = game.String()
	// very interesting https://stackoverflow.com/questions/46987513/handling-dangling-pointers-in-go
	(*db)[botID].lastPlayedMoveID = &playerID
	return nil
}

/*
assumes the potential returned board is valid
*/
func currentBoardHistory(botID int, db DBType) (string, error) {
	if _, ok := (*db)[botID]; !ok {
		// doesn't contain botID
		return "", fmt.Errorf("bot id doesn't exist")
	}
	// check if a game is running
	if game := (*db)[botID]; game == nil {
		return "", fmt.Errorf("no game running")
	} else {
		return game.gameState, nil
	}
}

func endGame(botID int, db DBType) error {
	if _, ok := (*db)[botID]; !ok {
		// doesn't contain botID
		return fmt.Errorf("bot id doesn't exist")
	}
	// check if a game is running
	if (*db)[botID] == nil {
		return fmt.Errorf("no game running")
	}
	(*db)[botID] = nil
	return nil
}

func main() {

	println("Server init!")
	server := gin.Default()
	db := map[int]*Game{}
	// API endpoints

	server.GET("/", func(c *gin.Context) {
		c.String(200, "Wrong endpoint : )")
	})

	server.GET("/register_client", func(c *gin.Context) {
		botID := c.Query("botID")
		res, err := strconv.ParseInt(botID, 10, 32)
		if err == nil {
			err = registerClient(int(res), &db)
			if err != nil {
				c.String(400, "%v", err)
			}
		} else {
			c.String(400, "invalid bot id: %v got: %v", err, botID)
		}
	})

	server.GET("/start_game", func(c *gin.Context) {
		botID := c.Query("botID")
		res, err := strconv.ParseInt(botID, 10, 32)
		if err == nil {
			err = startGame(int(res), &db)
			if err != nil {
				c.String(400, "%v", err)
			}
		} else {
			c.String(400, "invalid bot id: %v", err)
		}
	})

	server.GET("/play_move", func(c *gin.Context) {
		// trying a different error handling approach
		botID := c.Query("botID")
		playerID := c.Query("playerID")
		move_string := c.Query("move")
		var move_err error = nil
		botID_int, botID_err := strconv.ParseInt(botID, 10, 32)
		playerID_int, playerID_err := strconv.ParseInt(playerID, 10, 32)
		if move_string == "" {
			move_err = fmt.Errorf("invalid move string %v", move_string)
		}
		if botID_err == nil && playerID_err == nil && move_err == nil {
			// next step
			err := playMove(int(botID_int), int(playerID_int), move_string, &db)
			if err != nil {
				c.String(400, "%v", err)
			}
		}
		if botID_err != nil {
			c.String(400, "invalid bot id type: %v", botID_err)
		}
		if playerID_err != nil {
			c.String(400, "invalid player id type: %v", playerID_err)
		}
		if move_err != nil {
			c.String(400, "invalid move: %v", move_err)
		}
	})

	server.GET("/current_board_history", func(c *gin.Context) {
		botID := c.Query("botID")
		res, err := strconv.ParseInt(botID, 10, 32)
		if err == nil {
			boardState, err := currentBoardHistory(int(res), &db)
			if err != nil {
				c.String(400, "%v", err)
			} else {
				c.String(200, boardState)
			}
		} else {
			c.String(400, "invalid bot id: %v", err)
		}
	})

	server.GET("/end_game", func(c *gin.Context) {
		botID := c.Query("botID")
		res, err := strconv.ParseInt(botID, 10, 32)
		if err == nil {
			err = endGame(int(res), &db)
			if err != nil {
				c.String(400, "%v", err)
			}
		} else {
			c.String(400, "invalid bot id: %v", err)
		}
	})

	server.Run()
}
