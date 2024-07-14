package main

import (
	"fmt"
)

func main() {
	fmt.Println("This is a Go program compiled to WebAssembly.")
	fmt.Println("Seting up the Page")
	setupPage()
	currentBoard = createBoard(5, 5)
	currentBoard.print()
	displayBoard(currentBoard)

	select {}
}
