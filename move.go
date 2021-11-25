package main

type moveResult interface {
}

type moveTouchResult struct {
	gopher int
}

type moveBlockedResult struct {
}

type moveSuccessResult struct {
}

type movePushResult struct {
}

type moveEatResult struct {
	score int
}
