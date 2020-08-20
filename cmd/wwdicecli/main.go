package main

import (
	"flag"
	"fmt"
	"github.com/klaital/wwdice/pkg/dice"
	"math/rand"
	"os"
	"time"
)

func main() {
	var diceStr string
	flag.StringVar(&diceStr, "d", "1", "A wwdice formatstring. Takes the form of 5d6, to roll 5 10-sided dice against difficulty 6. Append a '!', to explode any 10's")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())
	d, err := dice.ParseDiceString(diceStr)
	if err != nil {
		fmt.Printf("Failed to parse input: %v", err)
		os.Exit(1)
	}

	rolls := d.RollDice()
	results := d.CalculateResults(rolls)

	if results.Botch {
		fmt.Printf("BOTCH!\t%v\n", rolls)
	} else {
		fmt.Printf("%d successes\t%v", results.Successes, rolls)
	}
}
