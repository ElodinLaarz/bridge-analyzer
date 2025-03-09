package main

import (
	"fmt"

	"github.com/ElodinLaarz/bridge-analyzer"
)

func main() {
	fmt.Println("Welcome to Bridge Analyzer!")

	bridgeGame := game.New()
	for i := 0; i < 2; i++ {
		fmt.Println("\n\nDealing a new game...")
		bridgeGame.Deal()
		fmt.Print(bridgeGame)
	}
}
