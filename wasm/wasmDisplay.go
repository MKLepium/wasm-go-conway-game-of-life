package main

import (
	"fmt"
	"strconv"
	"syscall/js"
	"time"
)

var (
	selectedCellX int = -1
	selectedCellY int = -1
)

var ticker *time.Ticker
var stopChan chan bool

// Global performance tracking
var totalDuration time.Duration
var updateCount int

func setupPage() {
	doc := js.Global().Get("document")
	body := doc.Get("body")
	body.Set("innerHTML", "")
	body.Set("style", "background-color: darkgrey; color: white;")

	setupHeader()
	setupStyle()
	setupFooter()
	createControlButtons()
	createButtonCreateNewBoard()
}

func createStartButton(doc js.Value, autoUpdateContainer js.Value, intervalInput js.Value) {
	startButton := doc.Call("createElement", "button")
	startButton.Set("innerText", "Start Auto Update")
	startButton.Set("style", "margin: 10px; padding: 10px; font-size: 18px; font-family: sans-serif; background-color: green; color: white; border: 1px solid black; border-radius: 5px;")
	startButtonClick := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Println("Auto Update Started")
		intervalStr := intervalInput.Get("value").String()
		interval, err := strconv.Atoi(intervalStr)
		if err != nil || interval < 60 {
			interval = 60
			// Minimum interval lower than that can make the browser unresponsive
		}
		ticker = time.NewTicker(time.Duration(interval) * time.Millisecond)
		stopChan = make(chan bool)
		totalDuration = 0
		updateCount = 0

		go func() {
			for {
				select {
				case <-ticker.C:
					startTime := time.Now()
					js.Global().Call("requestAnimationFrame", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
						updateBoard()
						duration := time.Since(startTime)
						totalDuration += duration
						updateCount++
						fmt.Printf("Update took: %v. Average time: %v\n", duration, totalDuration/time.Duration(updateCount))
						return nil
					}))
				case <-stopChan:
					return
				}
			}
		}()
		return nil
	})
	startButton.Call("addEventListener", "click", startButtonClick)
	autoUpdateContainer.Call("appendChild", startButton)
}

func createStopButton(doc js.Value, autoUpdateContainer js.Value) {
	stopButton := doc.Call("createElement", "button")
	stopButton.Set("innerText", "Stop Auto Update")
	stopButton.Set("style", "margin: 10px; padding: 10px; font-size: 18px; font-family: sans-serif; background-color: red; color: white; border: 1px solid black; border-radius: 5px;")
	stopButtonClick := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if ticker != nil {
			ticker.Stop()
			stopChan <- true
			fmt.Println("Auto Update Stopped")
		}
		return nil
	})
	stopButton.Call("addEventListener", "click", stopButtonClick)
	autoUpdateContainer.Call("appendChild", stopButton)
}

func createNextButton(doc js.Value, autoUpdateContainer js.Value) {
	nextButton := doc.Call("createElement", "button")
	nextButton.Set("innerText", "Next")
	nextButton.Set("style", "margin: 10px; padding: 10px; font-size: 18px; font-family: sans-serif; background-color: blue; color: white; border: 1px solid black; border-radius: 5px;")
	nextButtonClick := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		startTime := time.Now()
		updateBoard()
		duration := time.Since(startTime)
		fmt.Printf("Single update took: %v\n", duration)
		return nil
	})
	nextButton.Call("addEventListener", "click", nextButtonClick)
	autoUpdateContainer.Call("appendChild", nextButton)
}

func createIntervalInput(doc js.Value, autoUpdateContainer js.Value) js.Value {
	intervalInput := doc.Call("createElement", "input")
	intervalInput.Set("type", "number")
	intervalInput.Set("placeholder", "Interval in milliseconds")
	intervalInput.Set("style", "margin: 10px; padding: 10px; font-size: 18px; font-family: sans-serif; background-color: white; color: black; border: 1px solid black; border-radius: 5px;")
	autoUpdateContainer.Call("appendChild", intervalInput)
	return intervalInput
}

func createControlButtons() {
	doc := js.Global().Get("document")
	body := doc.Get("body")

	autoUpdateContainer := doc.Call("createElement", "div")
	autoUpdateContainer.Set("style", "margin-bottom: 20px;")

	intervalInput := createIntervalInput(doc, autoUpdateContainer)
	createStartButton(doc, autoUpdateContainer, intervalInput)
	createStopButton(doc, autoUpdateContainer)
	createNextButton(doc, autoUpdateContainer)

	body.Call("appendChild", autoUpdateContainer)
}

func createButtonCreateNewBoard() {
	doc := js.Global().Get("document")
	body := doc.Get("body")

	// Container for Board Creation controls
	boardCreationContainer := doc.Call("createElement", "div")

	// create two input fields for width and height
	widthInput := doc.Call("createElement", "input")
	widthInput.Set("type", "number")
	widthInput.Set("placeholder", "Width")
	widthInput.Set("style", "margin: 10px; padding: 10px; font-size: 18px; font-family: sans-serif; background-color: white; color: black; border: 1px solid black; border-radius: 5px;")
	boardCreationContainer.Call("appendChild", widthInput)

	heightInput := doc.Call("createElement", "input")
	heightInput.Set("type", "number")
	heightInput.Set("placeholder", "Height")
	heightInput.Set("style", "margin: 10px; padding: 10px; font-size: 18px; font-family: sans-serif; background-color: white; color: black; border: 1px solid black; border-radius: 5px;")
	boardCreationContainer.Call("appendChild", heightInput)

	button := doc.Call("createElement", "button")
	button.Set("innerText", "Create New Board")
	button.Set("style", "margin: 10px; padding: 10px; font-size: 18px; font-family: sans-serif; background-color: white; color: black; border: 1px solid black; border-radius: 5px;")
	buttonClick := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		widthStr := widthInput.Get("value").String()   // Get the width as string
		heightStr := heightInput.Get("value").String() // Get the height as string

		width, err := strconv.Atoi(widthStr) // Convert width to int
		if err != nil {
			fmt.Println("Error converting width to integer:", err)
			return nil
		}

		height, err := strconv.Atoi(heightStr) // Convert height to int
		if err != nil {
			fmt.Println("Error converting height to integer:", err)
			return nil
		}

		fmt.Println("Create New Board Button Clicked")
		currentBoard = createBoard(width, height)
		displayBoard(currentBoard)
		return nil
	})
	button.Call("addEventListener", "click", buttonClick)
	boardCreationContainer.Call("appendChild", button)

	// Append the Board Creation container to the body below the Auto Update container
	body.Call("appendChild", boardCreationContainer)
}

func getOrCreateTable() js.Value {
	doc := js.Global().Get("document")
	table := doc.Call("getElementById", "gameTable")
	if table.IsNull() {
		table = doc.Call("createElement", "Table")
		table.Set("id", "gameTable")
		table.Set("align", "center")
		table.Get("style").Set("width", "100%")
		table.Get("style").Set("border", "1px solid black")
		table.Get("style").Set("height", "100%")
		table.Get("style").Set("borderCollapse", "collapse")
		table.Get("style").Set("tableLayout", "fixed")
		table.Get("style").Set("margin", "auto")
	}
	return table
}

func displayBoard(b Board) {
	doc := js.Global().Get("document")
	body := doc.Get("body")
	table := getOrCreateTable()
	clearTable(table)

	for x, row := range b.cells {
		tr := doc.Call("createElement", "tr")
		for y := range row {
			td := doc.Call("createElement", "td")
			if b.isCellAlive(y, x) {
				td.Get("classList").Call("add", "alive")
			} else {
				td.Get("classList").Call("add", "dead")
			}
			td = addClickListener(td, x, y)

			tr.Call("appendChild", td)
		}
		table.Call("appendChild", tr)
	}
	body.Call("appendChild", table)
	fmt.Println("Displaying board...")
}

func clearTable(table js.Value) {
	for !table.Get("firstChild").IsNull() {
		table.Call("removeChild", table.Get("firstChild"))
	}
}

func createFooterButtons(footer js.Value) {
	doc := js.Global().Get("document")

	// Button to set selected cell alive
	buttonAlive := doc.Call("createElement", "button")
	buttonAlive.Set("innerText", "Set Alive")
	buttonAlive.Set("style", "margin: 10px; padding: 10px; font-size: 18px; background-color: white; color: black; border: 1px solid black; border-radius: 5px;")
	buttonAliveClick := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if selectedCellX != -1 && selectedCellY != -1 {
			fmt.Println("Set Alive clicked for selected cell")
			currentBoard.setAlive(selectedCellY, selectedCellX)
			displayBoard(currentBoard)
			updateFooterText()
		}
		return nil
	})
	buttonAlive.Call("addEventListener", "click", buttonAliveClick)

	// Button to set selected cell dead
	buttonDead := doc.Call("createElement", "button")
	buttonDead.Set("innerText", "Set Dead")
	buttonDead.Set("style", "margin: 10px; padding: 10px; font-size: 18px; background-color: white; color: black; border: 1px solid black; border-radius: 5px;")
	buttonDeadClick := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if selectedCellX != -1 && selectedCellY != -1 {
			fmt.Println("Set Dead clicked for selected cell")
			currentBoard.setDead(selectedCellY, selectedCellX)
			displayBoard(currentBoard)
			updateFooterText()
		}
		return nil
	})
	buttonDead.Call("addEventListener", "click", buttonDeadClick)

	footer.Call("appendChild", buttonAlive)
	footer.Call("appendChild", buttonDead)
}

func addClickListener(td js.Value, x, y int) js.Value {
	tdClick := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		selectedCellX = x
		selectedCellY = y
		updateFooterText()
		return nil
	})
	td.Call("addEventListener", "click", tdClick)
	return td
}

func updateFooterText() {
	if selectedCellX == -1 || selectedCellY == -1 {
		js.Global().Get("document").Call("getElementById", "footerText").Set("innerText", "No cell selected")
	} else {
		js.Global().Get("document").Call("getElementById", "footerText").Set("innerText", fmt.Sprintf("Cell (%d, %d) selected", selectedCellX, selectedCellY))
	}
}

func setupFooter() {
	doc := js.Global().Get("document")
	footer := doc.Call("createElement", "footer")
	footer.Set("id", "footer")
	footer.Get("style").Set("position", "fixed")
	footer.Get("style").Set("bottom", "0")
	footer.Get("style").Set("width", "100%")
	footer.Get("style").Set("color", "white")
	footer.Get("style").Set("textAlign", "center")
	footer.Get("style").Set("backgroundColor", "#808080")
	footer.Get("style").Set("padding", "10px")
	footer.Get("style").Set("fontSize", "18px")
	footer.Get("style").Set("fontFamily", "sans-serif")
	doc.Get("body").Call("appendChild", footer)

	textContainer := doc.Call("createElement", "div")
	textContainer.Set("id", "footerText")
	textContainer.Set("innerText", "No cell selected")
	footer.Call("appendChild", textContainer)

	createFooterButtons(footer)
}

func setupHeader() {
	doc := js.Global().Get("document")
	body := doc.Get("body")

	header := doc.Call("createElement", "HEADER")
	header.Set("innerText", "Conways Game of Life using Golang and WebAssembly")
	header.Get("style").Set("textAlign", "center")
	header.Get("style").Set("fontSize", "24px")
	header.Get("style").Set("fontFamily", "sans-serif")
	header.Get("style").Set("width", "100%")
	header.Get("style").Set("padding", "20px")
	createDropdownMenu(header)

	body.Call("appendChild", header)

}

func createDropdownMenu(header js.Value) {
	doc := js.Global().Get("document")

	dropdown := doc.Call("createElement", "select")
	dropdown.Set("id", "boardOptions")
	dropdown.Set("style", "margin: 10px; padding: 10px; font-size: 18px; font-family: sans-serif; border: 1px solid black; border-radius: 5px;")

	options := map[string]Board{
		"Cannon (50x50)": GosperGun(50, 50),
		"Empty (50x50)":  createBoard(50, 50),
	}

	for name := range options {
		option := doc.Call("createElement", "option")
		option.Set("value", name)
		option.Set("innerText", name)
		dropdown.Call("appendChild", option)
	}

	button := doc.Call("createElement", "button")
	button.Set("innerText", "Create Selected Board")
	button.Set("style", "margin: 10px; padding: 10px; font-size: 18px; font-family: sans-serif; background-color: white; color: black; border: 1px solid black; border-radius: 5px;")
	buttonClick := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		selectedOption := dropdown.Get("value").String()
		currentBoard = options[selectedOption]
		displayBoard(currentBoard)
		return nil
	})
	button.Call("addEventListener", "click", buttonClick)

	header.Call("appendChild", dropdown)
	header.Call("appendChild", button)
}

func setupStyle() {
	doc := js.Global().Get("document")
	style := doc.Call("createElement", "style")
	css := `
		.alive {
			background-color: green;  /* Green background for alive cells */
		}
		.dead {
			background-color: gray;   /* Gray background for dead cells */
		}
		#gameTable td {
			width: 6px;
			height: 6px;
			border: 1px solid black;  /* Adds a border to each cell */
		}
	`
	style.Set("innerHTML", css)
	doc.Get("head").Call("appendChild", style)
}
