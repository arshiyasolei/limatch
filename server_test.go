package main

import (
	"reflect"
	"strings"
	"testing"

	"github.com/notnil/chess"
)

func TestRegisterClient(t *testing.T) {
	testDB := map[int]*Game{}
	registerClient(12345, &testDB)
	want := map[int]*Game{12345: nil}
	if !reflect.DeepEqual(testDB, want) {
		t.Errorf("Got: %v Want: %v ", testDB, want)
	}
}

func TestStartGame(t *testing.T) {
	testDB := map[int]*Game{12345: nil}
	startGame(12345, &testDB)
	want := map[int]*Game{12345: &Game{}}
	if !reflect.DeepEqual(testDB, want) {
		t.Errorf("Got: %v Want: %v ", testDB, want)
	}

	// game should not start because a game is already running
	err := startGame(12345, &testDB)
	if err == nil {
		t.Error("Expected game to not start")
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
	testDB := map[int]*Game{12345: &Game{new(int), ""}}
	want := game.String()

	err = playMove(12345, 1, "e4", &testDB)
	if err != nil {
		panic(err)
	}
	if testDB[12345].gameState != want {
		t.Errorf("want: %v got: %v", want, testDB[12345].gameState)
	}

	err = playMove(12345, 1, "e5", &testDB)
	if err == nil {
		t.Errorf("Wanted error for invalid player id but got no errors")
	}
}

func TestCurrentBoardHistory(t *testing.T) {
	currentPGN := `totally valid history`
	testDB := map[int]*Game{12345: &Game{new(int), currentPGN}}
	want, _ := currentBoardHistory(12345, &testDB)
	if testDB[12345].gameState != want {
		t.Errorf("want: %v got: %v", want, testDB[12345].gameState)
	}
}

func TestEndGame(t *testing.T) {
	testDB := map[int]*Game{12345: &Game{new(int), ""}}
	endGame(12345, &testDB)
	want := map[int]*Game{12345: nil}
	if !reflect.DeepEqual(testDB, want) {
		t.Errorf("Got: %v Want: %v ", testDB, want)
	}
}
