package dice

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strings"
)

type DiceConfig struct {
	Count uint8
	Difficulty uint8
	Explode bool
}

var ErrNoDice = errors.New("invalid dice number")
var ErrBadDifficulty = errors.New("invalid difficulty number")
var ErrFailedToParse = errors.New("failed to parse dice string")

// ParseDiceString parses a useful command out of a string in the format XdY!
// X = number of d10's to roll
// Y = difficulty of the roll. Set to 0 to
func ParseDiceString(s string) (DiceConfig, error) {
	var dice DiceConfig
	n, err := fmt.Sscanf(s, "%dd%d", &dice.Count, &dice.Difficulty)
	if err != nil {
		return dice, ErrFailedToParse
	}
	if n != 2 {
		return dice, ErrFailedToParse
	}

	dice.Explode = strings.Contains(s, "!")

	// A difficulty of greater than 10 is not allowed, as we're rolling d10's.
	if dice.Difficulty > 10 || dice.Difficulty == 0 {
		return dice, ErrBadDifficulty
	}

	// Can't roll 0 dice, even though it parses correctly
	if dice.Count == 0 {
		return dice, ErrNoDice
	}

	return dice, nil
}

type RollResults struct {
	Successes uint8
	Failures uint8
	Botch bool
}
// CalculateResults takes in the results of a dice roll and tells you how many successes and failures you got.
func (d DiceConfig) CalculateResults(results []uint8) RollResults {
	var successes uint8
	var failures uint8
	botch := false
	for _, roll := range results {
		if roll >= d.Difficulty {
			successes += 1
		} else if roll == 1 {
			failures += 1
		}
	}
	if failures > 0 && successes == 0 {
		botch = true
	} else if failures > successes {
		successes = 0
	} else {
		successes -= failures
	}

	return RollResults{
		Successes: successes,
		Failures:  failures,
		Botch:     botch,
	}
}

// RollDice generates a set of random dice rolls according to the parsed config
func (d DiceConfig) RollDice() []uint8 {
	diceToRoll := d.Count
	results := make([]uint8, 0, d.Count)
	for len(results) < int(diceToRoll) {
		// Get a number between 1 and 10, inclusive
		i := (rand.Uint32() % 10) + 1
		results = append(results, uint8(i))
		if i == 10 && d.Explode {
			diceToRoll += 1
		}
	}

	// Sort them from highest to lowest
	sort.Sort(UInt8Slice(results))
	return results
}

// UInt8Slice implements the interface needed to use the sort package.
type UInt8Slice []uint8
func (p UInt8Slice) Len() int           { return len(p) }
func (p UInt8Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p UInt8Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
