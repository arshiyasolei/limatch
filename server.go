package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/notnil/chess"
)

type Game struct {
	lastPlayedMovePlayerId *int   // optional
	gameState              string // PGN or something
}

type DBType *map[int]*Game

/*
In memory 'database'
*/
func registerClient(bot_id int, db DBType) error {
	if _, ok := (*db)[bot_id]; ok {
		// contains bot_id already
		return fmt.Errorf("contains bot_id already")
	} else {
		(*db)[bot_id] = nil
		return nil
	}
}

func startGame(bot_id int, db DBType) error {
	if _, ok := (*db)[bot_id]; !ok {
		// doesn't contain bot_id
		return fmt.Errorf("bot id doesn't exist")
	}
	// check if a another game is running
	if game := (*db)[bot_id]; game != nil {
		return fmt.Errorf("another game is running already")
	}
	(*db)[bot_id] = &Game{nil, ""}
	return nil
}

func playMove(bot_id int, player_id int, move string, db DBType) error {
	if _, ok := (*db)[bot_id]; !ok {
		// doesn't contain bot_id
		return fmt.Errorf("bot id doesn't exist")
	}
	// check if a game is running
	if game := (*db)[bot_id]; game == nil {
		return fmt.Errorf("no game running")
	}

	// check if we are not making the a new move with the same player
	// that played the previous move
	if game := (*db)[bot_id]; game.lastPlayedMovePlayerId != nil && *game.lastPlayedMovePlayerId == player_id {
		return fmt.Errorf("consecutive moves detected")
	}

	// instantiate a board and load current state
	pgn, err := chess.PGN(strings.NewReader((*db)[bot_id].gameState))
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
	(*db)[bot_id].gameState = game.String()
	// very interesting https://stackoverflow.com/questions/46987513/handling-dangling-pointers-in-go
	(*db)[bot_id].lastPlayedMovePlayerId = &player_id
	return nil
}

/*
assumes the potential returned board is valid
*/
func currentBoardHistory(bot_id int, db DBType) (string, error) {
	if _, ok := (*db)[bot_id]; !ok {
		// doesn't contain bot_id
		return "", fmt.Errorf("bot id doesn't exist")
	}
	// check if a game is running
	if game := (*db)[bot_id]; game == nil {
		return "", fmt.Errorf("no game running")
	} else {
		return game.gameState, nil
	}
}

func endGame(bot_id int, db DBType) error {
	if _, ok := (*db)[bot_id]; !ok {
		// doesn't contain bot_id
		return fmt.Errorf("bot id doesn't exist")
	}
	// check if a game is running
	if (*db)[bot_id] == nil {
		return fmt.Errorf("no game running")
	}
	(*db)[bot_id] = nil
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
		bot_id := c.Query("bot_id")
		res, err := strconv.ParseInt(bot_id, 10, 32)
		if err == nil {
			err = registerClient(int(res), &db)
			if err != nil {
				c.String(400, "%v", err)
			}
		} else {
			c.String(400, "invalid bot id: %v got: %v", err, bot_id)
		}
	})

	server.GET("/start_game", func(c *gin.Context) {
		bot_id := c.Query("bot_id")
		res, err := strconv.ParseInt(bot_id, 10, 32)
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
		bot_id := c.Query("bot_id")
		player_id := c.Query("player_id")
		move_string := c.Query("move")
		var move_err error = nil
		bot_id_int, bot_id_err := strconv.ParseInt(bot_id, 10, 32)
		player_id_int, player_id_err := strconv.ParseInt(player_id, 10, 32)
		if move_string == "" {
			move_err = fmt.Errorf("invalid move string %v", move_string)
		}
		if bot_id_err == nil && player_id_err == nil && move_err == nil {
			// next step
			err := playMove(int(bot_id_int), int(player_id_int), move_string, &db)
			if err != nil {
				c.String(400, "%v", err)
			}
		}
		if bot_id_err != nil {
			c.String(400, "invalid bot id type: %v", bot_id_err)
		}
		if player_id_err != nil {
			c.String(400, "invalid player id type: %v", player_id_err)
		}
		if move_err != nil {
			c.String(400, "invalid move: %v", move_err)
		}
	})

	server.GET("/current_board_history", func(c *gin.Context) {
		bot_id := c.Query("bot_id")
		res, err := strconv.ParseInt(bot_id, 10, 32)
		if err == nil {
			board_state, err := currentBoardHistory(int(res), &db)
			if err != nil {
				c.String(400, "%v", err)
			} else {
				c.String(200, board_state)
			}
		} else {
			c.String(400, "invalid bot id: %v", err)
		}
	})

	server.GET("/end_game", func(c *gin.Context) {
		bot_id := c.Query("bot_id")
		res, err := strconv.ParseInt(bot_id, 10, 32)
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
