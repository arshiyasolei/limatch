package main

import (
	"reflect"
	"strings"
	"testing"

	"github.com/notnil/chess"
)

func TestRegisterClient(t *testing.T) {
	test_db := map[int]*Game{}
	registerClient(12345, &test_db)
	want := map[int]*Game{12345: nil}
	if !reflect.DeepEqual(test_db, want) {
		t.Errorf("Got: %v Want: %v ", test_db, want)
	}
}

func TestStartGame(t *testing.T) {
	test_db := map[int]*Game{12345: nil}
	startGame(12345, &test_db)
	want := map[int]*Game{12345: &Game{}}
	if !reflect.DeepEqual(test_db, want) {
		t.Errorf("Got: %v Want: %v ", test_db, want)
	}
}

func TestPlayMove(t *testing.T) {

	testPGN := ``
	testPGNReader := strings.NewReader(testPGN)
	pgn, err := chess.PGN(testPGNReader)
	if err != nil {
		panic("failed to parse PGN with error" + err.Error())
	}
	game := chess.NewGame(pgn)
	game.MoveStr("e4")
	test_db := map[int]*Game{12345: &Game{new(int), ""}}
	want := game.String()

	err = playMove(12345, 1, "e4", &test_db)
	if err != nil {
		panic(err)
	}
	if test_db[12345].gameState != want {
		t.Errorf("want: %v got: %v", want, test_db[12345].gameState)
	}
}

func TestCurrentBoardHistory(t *testing.T) {
	currentPGN := `totally valid history`
	test_db := map[int]*Game{12345: &Game{new(int), currentPGN}}
	want, _ := currentBoardHistory(12345, &test_db)
	if test_db[12345].gameState != want {
		t.Errorf("want: %v got: %v", want, test_db[12345].gameState)
	}
}

func TestEndGame(t *testing.T) {
	test_db := map[int]*Game{12345: &Game{new(int), ""}}
	endGame(12345, &test_db)
	want := map[int]*Game{12345: nil}
	if !reflect.DeepEqual(test_db, want) {
		t.Errorf("Got: %v Want: %v ", test_db, want)
	}
}
