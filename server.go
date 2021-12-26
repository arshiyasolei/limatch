package main

import (
	"fmt"
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
		return nil
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
	if game := (*db)[bot_id]; game != nil && game.lastPlayedMovePlayerId != nil {
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
	if (*db)[bot_id].lastPlayedMovePlayerId == nil {
		return fmt.Errorf("no game running")
	}
	(*db)[bot_id] = nil
	return nil
}

func main() {
	chess.NewGame()
	fmt.Errorf("d")
	println("Server init!")
	server := gin.Default()

	// API endpoints

	server.GET("/", func(c *gin.Context) {
		c.String(200, "Wrong endpoint : )")
	})

	server.Run()
}
