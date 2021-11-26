package main

type moveResult interface {
}

type moveTouchResult struct {
	x, y   int
	gopher int
}

type moveBlockedResult struct {
	x, y int
}

type moveSuccessResult struct {
	x, y int
}

type movePushResult struct {
	x, y int
}

type moveEatResult struct {
	x, y  int
	score int
}
