package characters

import (
	"gopkg.in/yaml.v2"
	"strconv"
	"strings"
)

// ParseFormula takes in strings of the format
// attribute+ability[+1][!]
// And returns the score's tags, any extra dice added in, and whether to explode 10's.
func ParseFormula(formula string) (tags []string, bonuses []uint8, explode bool, err error) {
	// Init
	tags = make([]string, 0)
	bonuses = make([]uint8, 0)
	explode = strings.Contains(formula, "!")

	// Check for exploding 10's symbol
	if explode {
		formula = strings.Replace(formula, "!", "", 1)
	}

	// Parse out the variables and primitives to add up
	tokens := strings.Split(formula, "+")
	for i, token := range tokens {
		// Try to cast to an int
		intVal, err := strconv.Atoi(token)
		if err == nil && intVal > 0 {
			bonuses = append(bonuses, uint8(intVal))
		}

		// Treat it as a tag name
		tags = append(tags, tokens[i])
	}

	return tags, bonuses, explode, nil
}

func (c Character) CountDice(tags []string, bonuses []uint8) uint8 {
	var sum uint8
	for _, tag := range tags {
		sum += c[tag]
	}
	for _, bonus := range bonuses {
		sum += bonus
	}
	return sum
}

type Character map[string]uint8
func (c Character) Validate() bool {
	nonzeroFields := []string{"str", "dex", "sta", "man", "app", "cha", "per", "int", "wis"}
	for _, field := range nonzeroFields {
		if c[field] == 0 {
			return false
		}
	}
	return true
}
func (c Character) ToString() string {
	characterBytes, err := yaml.Marshal(c)
	if err != nil {
		return ""
	}
	return string(characterBytes)
}
