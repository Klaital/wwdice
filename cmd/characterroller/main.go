package main

import (
	"fmt"
	"github.com/klaital/wwdice/pkg/characters"
	"github.com/klaital/wwdice/pkg/dice"
)

func main() {
	c := characters.Character{
		Strength: 5,
		Athletics: 2,
	}

	tags, bonuses, explode, _ := characters.ParseFormula("str+athletics")
	diceCount := c.CountDice(tags, bonuses)
	fmt.Printf("Character has str+athletics = %d\n", diceCount)
	diceCfg := dice.DiceConfig{
		Count: diceCount,
		Explode: explode,
		Difficulty: 6,
	}
	rolls := diceCfg.RollDice()
	fmt.Printf("%v\n", rolls)
}
