package dice

import "testing"

func TestParseDiceString(t *testing.T) {
	checkDiceParsing(" 3d6 ",
		DiceConfig{
			Count:      3,
			Difficulty: 6,
			Explode:    false,
		},
		nil,
		t)

	checkDiceParsing("3d6",
		DiceConfig{
			Count:      3,
			Difficulty: 6,
			Explode:    false,
		},
		nil,
		t)

	checkDiceParsing("10d7!",
		DiceConfig{
			Count:      10,
			Difficulty: 7,
			Explode:    true,
		},
		nil,
		t)

	checkDiceParsing("3d11",
		DiceConfig{
			Count:      3,
			Difficulty: 11,
			Explode:    false,
		},
		ErrBadDifficulty,
		t)

	checkDiceParsing("10d0",
		DiceConfig{
			Count:      10,
			Difficulty: 0,
			Explode:    false,
		},
		ErrBadDifficulty,
		t)

	checkDiceParsing("0d6",
		DiceConfig{
			Count:      0,
			Difficulty: 6,
			Explode:    false,
		},
		ErrNoDice,
		t)

	//
	// These should all fail to parse
	//
	checkDiceParsing("xd6",
		DiceConfig{},
		ErrFailedToParse,
		t)
	checkDiceParsing("6dy",
		DiceConfig{},
		ErrFailedToParse,
		t)
}

func checkDiceParsing(original string, expected DiceConfig, expectedErr error, t *testing.T) {
	parsed, err := ParseDiceString(original)
	if err != expectedErr {
		t.Errorf("expected parsing error %v from %s. Got %v", expectedErr, original, err)
		return
	}
	// Don't bother testing the results if there was a parse error
	if expectedErr == nil {
		if expected.Count != parsed.Count {
			t.Errorf("Incorrect count from %s. Expected %d, got %d.", original, expected.Count, parsed.Count)
		}
		if expected.Difficulty != parsed.Difficulty {
			t.Errorf("Incorrect difficulty from %s. Expected %d, got %d.", original, expected.Difficulty, parsed.Difficulty)
		}
		if expected.Explode != parsed.Explode {
			t.Errorf("Incorrect explode setting from %s. Expected %t, got %t.", original, expected.Explode, parsed.Explode)
		}
	}
}

func TestDiceConfig_CalculateResults(t *testing.T) {

	rolls := []uint8{1, 4, 6, 7, 10}
	results := DiceConfig{Count: 5, Difficulty: 6}.CalculateResults(rolls)
	checkResultsCalculation(rolls, RollResults{2, 1, false}, results, t)

	rolls = []uint8{1, 2, 3, 4, 5}
	results = DiceConfig{Count: 5, Difficulty: 6}.CalculateResults(rolls)
	checkResultsCalculation(rolls, RollResults{0, 1, true}, results, t)
}

func checkResultsCalculation(rolls []uint8, expected, actual RollResults, t *testing.T) {
	if expected.Successes != actual.Successes {
		t.Errorf("Incorrect success count from %v. Expected %d, got %d", rolls, expected.Successes, actual.Successes)
	}
	if expected.Failures != actual.Failures {
		t.Errorf("Incorrect failure count from %v. Expected %d, got %d", rolls, expected.Failures, actual.Failures)
	}
	if expected.Botch != actual.Botch {
		t.Errorf("Incorrect botch calculation from %v. Expected %t, got %t", rolls, expected.Botch, actual.Botch)
	}
}

func TestDiceConfig_RollDice(t *testing.T) {
	iterations := 100000
	var d DiceConfig
	var rolls []uint8

	// Validate that we never explode dice when told not to, and that we never
	// get dice values outside the expected range of 1-10
	d = DiceConfig{Count: 5, Difficulty: 6, Explode: false}
	for i := 0; i < iterations; i++ {
		rolls = d.RollDice()
		if len(rolls) != int(d.Count) {
			t.Errorf("Got an incorrect count of results - expected no exploding dice. Expected %d, got %d", d.Count, len(rolls))
		}
		for _, dieValue := range rolls {
			if dieValue == 0 || dieValue > 10 {
				t.Errorf("Got an incorrect value for one of the dice: %d", dieValue)
			}
		}
	}

	// Validate that sometimes dice explode when told to
	d = DiceConfig{Count: 10, Difficulty: 6, Explode: true}
	everExploded := false
	for i := 0; i < iterations; i++ {
		rolls = d.RollDice()
		if len(rolls) > int(d.Count) {
			everExploded = true
			// Make sure there was a number of 10's equal to the number of extra rolls returned
			expectedExplodeCount := 0
			for _, dieValue := range rolls {
				if dieValue == 10 {
					expectedExplodeCount += 1
				}
			}
			if expectedExplodeCount != len(rolls) - int(d.Count) {
				t.Errorf("Incorrect explosion count. Roll returned %d extra dice, but there were %d 10's in the set.", len(rolls) - int(d.Count), expectedExplodeCount)
			}
		}
	}
	if !everExploded {
		t.Errorf("The roller never exploded any dice over %d iterations of rolling %d dice", iterations, d.Count)
	}

}