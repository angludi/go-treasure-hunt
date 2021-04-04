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

	ERR_DIRECTION = "You can not move to that direction!"
	GET_TREASURE  = "Yeaaay, you got the treasure!"
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
	title := "TREASURE HUNT"

	// Header
	DrawLine(len(title))
	fmt.Println(title)
	DrawLine(len(title))

	self.X = 1
	self.Y = 1

	DrawGrid()

	for {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Where do you want to go? ")

		key, _, err := reader.ReadRune()
		if err != nil {
			fmt.Println(err.Error())
		}

		if key == 'q' {
			ExitTheGame()
		}

		response := Move(key, self)

		if response.GotTheTreasure {
			fmt.Println(GET_TREASURE)
			ExitTheGame()
		} else if response.WrongDirection {
			fmt.Println(ERR_DIRECTION)
		}

		self = response.Position

		DrawGrid()
	}

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

func DrawGrid() {
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
					} else if CheckTreasure(pos) {
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

	response.Position = position
	response.WrongDirection = true

	if !wrongKey {
		if CheckTreasure(newPosition) {
			response.Position = newPosition
			response.WrongDirection = false
			response.GotTheTreasure = true
		} else if !CheckOutside(newPosition) && !CheckObstacle(newPosition) {
			response.Position = newPosition
			response.WrongDirection = false
		}
	}

	return response
}

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

func CheckOutside(position Position) bool {
	if position.X < 1 || position.X > GRID_WIDTH || position.Y < 1 || position.Y > GRID_HEIGHT {
		return true
	}
	return false
}

func CheckTreasure(position Position) bool {
	if position.X == treasure.X && position.Y == treasure.Y {
		return true
	}
	return false
}

func ExitTheGame() {
	fmt.Println("Thank's for playing!")
	os.Exit(0)
}
