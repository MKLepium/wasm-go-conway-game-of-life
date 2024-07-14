package main

import (
	"fmt"
)

type Cell struct {
	alive bool
}

type Board struct {
	cells [][]Cell
}

var currentBoard Board

func createBoard(width, height int) Board {
	board := Board{cells: make([][]Cell, height)}
	for i := range board.cells {
		board.cells[i] = make([]Cell, width)
	}
	return board
}

func updateBoard() {
	currentBoard = currentBoard.nextGeneration()
	displayBoard(currentBoard)
}

func (b *Board) setCell(x, y int, alive bool) {
	b.cells[y][x].alive = alive
}

func (b *Board) isCellAlive(x, y int) bool {
	return b.cells[y][x].alive
}

func (b *Board) toggleCell(x, y int) {
	b.setCell(x, y, !b.isCellAlive(x, y))
}

func (b *Board) setAlive(x, y int) {
	b.setCell(x, y, true)
}

func (b *Board) setDead(x, y int) {
	b.setCell(x, y, false)
}

func (b *Board) print() {
	for x, row := range b.cells {
		for y := range row {
			if b.isCellAlive(y, x) {
				fmt.Print("X")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func (b *Board) countNeighbors(x, y int) int {
	count := 0
	maxX := len(b.cells[0]) - 1 // Maximum X index
	maxY := len(b.cells) - 1    // Maximum Y index

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue // Skip the current cell itself
			}

			neighborX := x + j
			neighborY := y + i

			// Check if the neighbor's coordinates are within the board's bounds
			if neighborX >= 0 && neighborX <= maxX && neighborY >= 0 && neighborY <= maxY {
				if b.isCellAlive(neighborX, neighborY) {
					count++
				}
			}
		}
	}
	return count
}

func (b *Board) nextGeneration() Board {
	next := createBoard(len(b.cells[0]), len(b.cells))
	for x, row := range b.cells {
		for y := range row {
			neighbors := b.countNeighbors(y, x)
			alive := b.isCellAlive(y, x)
			if alive && (neighbors < 2 || neighbors > 3) {
				next.setCell(y, x, false)
			} else if !alive && neighbors == 3 {
				next.setCell(y, x, true)
			} else {
				next.setCell(y, x, alive)
			}
		}
	}
	return next
}
