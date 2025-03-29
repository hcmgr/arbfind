package main

import "fmt"

// For the given sport's matches, tests that the team
// positions are the same across bookies for that match.
//
// Example:
//
// Take the match 'Bulldogs v Magpies'.
// If, for the first bookie: outcomes[0] == 'Bulldogs' AND outcomes[1] == 'Magpies',
// we want to ensure this is true for the rest of the bookies.
// In general, this ensures team A and team B are always in the same position
// for each bookie in a match.
func testTeamsSamePositionSport(matches []Match) {
	for _, match := range matches {
		bookmakers := match.Bookmakers
		if len(bookmakers) == 0 || len(bookmakers[0].Markets) == 0 {
			continue
		}
		firstName := bookmakers[0].Markets[0].Outcomes[0].Name
		secondName := bookmakers[0].Markets[0].Outcomes[1].Name
		for i := 1; i < len(bookmakers); i++ {
			otherFirstName := bookmakers[i].Markets[0].Outcomes[0].Name
			otherSecondName := bookmakers[i].Markets[0].Outcomes[1].Name
			if firstName != otherFirstName || secondName != otherSecondName {
				fmt.Println("DISCREPENCY")
				return
			}
		}
	}
	fmt.Println("We good")
}
