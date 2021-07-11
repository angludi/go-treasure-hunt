package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	GRID_WIDTH  = 6
	GRID_HEIGHT = 4

	ERR_DIRECTION = "Oh no, you hit obstacle! Better luck next time!"
	GET_TREASURE  = "Yeaaay, you got the treasure!"
	NO_LUCK = "Better luck next time!"

	ERR_MOVE = "Invalid move. Must be numeric!"
)

type Position struct {
	X int
	Y int
}

type MoveResponse struct {
	Position       Position
	WrongDirection bool
	GotTheTreasure bool
}

var obstacles = [6][2]int{{2, 1}, {4, 2}, {6, 2}, {2, 3}, {3, 3}, {4, 3}}
var treasure = HideTheTreasure()
var self Position


func main() {
	var up, right, down int

	title := "TREASURE HUNT"

	// Header
	DrawLine(len(title))
	fmt.Println(title)
	DrawLine(len(title))

	// Set player location start point
	self.X = 1
	self.Y = 1

	DrawGrid(false)

	fmt.Println("Please input steps (numeric) for up, right and down (Separate by space):")
	var r = bufio.NewReader(os.Stdin)
	fmt.Fscanf(r, "%d %d %d", &up, &right, &down)

	response := Moves(up, right, down, self)

	var msgResult string
	if response.GotTheTreasure {
		msgResult = GET_TREASURE
	}else if response.WrongDirection {
		msgResult = ERR_DIRECTION
	} else {
		msgResult = NO_LUCK
	}

	self = response.Position

	DrawGrid(true)
	fmt.Println(msgResult)
	fmt.Println("Player Position:", self)
	fmt.Println("Treasure Position:", treasure)
	ExitTheGame()
}

func HideTheTreasure() Position {
	var save bool
	var location Position

	rand.Seed(time.Now().UnixNano())

	for !save {
		x := rand.Intn(GRID_WIDTH) + 1
		y := rand.Intn(GRID_HEIGHT) + 1

		// Hide the treasure in clear path, and not at player location
		if !CheckObstacle(Position{X: x, Y: y}) && (y != 1 && x != 1) {
			location.X = x
			location.Y = y
			save = true
		}
	}

	return location
}

func DrawLine(size int) {
	fmt.Println(strings.Repeat("-", size))
}

func DrawGrid(showTreasure bool) {
	// I'll give you a space
	fmt.Println()
	for i := GRID_HEIGHT + 1; i >= 0; i-- {
		for j := 0; j <= GRID_WIDTH+1; j++ {
			pos := Position{X: j, Y: i}
			// Top & Bottom Border
			if j == GRID_WIDTH+1 || i == GRID_HEIGHT+1 {
				fmt.Print("#")
				if j == GRID_WIDTH+1 {
					fmt.Println("")
				}
			} else {
				if i > 0 && j > 0 {
					if pos == self {
						// Player
						fmt.Print("X")
					} else if CheckObstacle(pos) {
						// Obstacle
						fmt.Print("#")
					} else if CheckTreasure(pos) && showTreasure {
						// The Treasure
						fmt.Print("$")
					} else {
						// Clear Path
						fmt.Print(".")
					}
				} else {
					// Left & Right Border
					fmt.Print("#")
				}
			}
		}
	}
}

func Moves(up, right, down int, position Position) (response MoveResponse) {
	for i := 0; i < up; i++ {
		response = Move('u', position)
		if response.WrongDirection {
			return
		}
		position = response.Position
	}

	for i := 0; i < right; i++ {
		response = Move('r', position)
		if response.WrongDirection {
			return
		}
		position = response.Position
	}

	for i := 0; i < down; i++ {
		response = Move('d', position)
		if response.WrongDirection {
			return
		}
		position = response.Position
	}

	return
}

// Move player by key pressed
func Move(key rune, position Position) (response MoveResponse) {
	var wrongKey bool
	var newPosition Position

	newPosition = position

	switch key {
	case 'l':
		newPosition.X = position.X - 1
		break
	case 'r':
		newPosition.X = position.X + 1
		break
	case 'u':
		newPosition.Y = position.Y + 1
		break
	case 'd':
		newPosition.Y = position.Y - 1
		break
	default:
		wrongKey = true
	}

	response.Position = newPosition
	response.WrongDirection = false

	if CheckTreasure(newPosition) {
		response.GotTheTreasure = true
	} else if CheckOutOfGrid(newPosition) || CheckObstacle(newPosition) || wrongKey {
		response.Position = position
		response.WrongDirection = true
	}

	return response
}

// Check  if posisition is an obstacle
func CheckObstacle(position Position) bool {
	for _, obstacle := range obstacles {
		x := obstacle[0]
		y := obstacle[1]

		if position.X == x && position.Y == y {
			return true
		}
	}
	return false
}

// Check if position is out of grid
func CheckOutOfGrid(position Position) bool {
	if position.X < 1 || position.X > GRID_WIDTH || position.Y < 1 || position.Y > GRID_HEIGHT {
		return true
	}
	return false
}

// Check if position is location of the treasure
func CheckTreasure(position Position) bool {
	if position.X == treasure.X && position.Y == treasure.Y {
		return true
	}
	return false
}

// Quit the game
func ExitTheGame() {
	fmt.Println("Thank's for playing!")
	os.Exit(0)
}
