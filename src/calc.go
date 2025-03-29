package main

import "fmt"

func checkArb(match *Match, outcomes []*Outcome, arbs *[]Arb) {
	var r float64 = 0
	for _, o := range outcomes {
		r += (1 / o.Price)
	}

	if r < 1 {
		*arbs = append(*arbs, constructArb(match, outcomes, r))
	}
}

func constructArb(match *Match, outcomes []*Outcome, R float64) Arb {
	return Arb{
		MatchId:      match.MatchId,
		SportKey:     match.SportKey,
		SportTitle:   match.SportTitle,
		CommenceTime: match.CommenceTime,
		HomeTeam:     match.HomeTeam,
		AwayTeam:     match.AwayTeam,
		Outcomes:     outcomes,
		R:            R,
	}
}

func findTwoWayMatchArbs(match *Match, arbs *[]Arb) {
	for i := 0; i < len(match.Bookmakers); i++ {
		bookmakerA := &match.Bookmakers[i]
		outcomesA := bookmakerA.Markets[0].Outcomes

		// prune non-two-way bookie
		if len(outcomesA) != 2 {
			continue
		}

		oA1, oA2 := &outcomesA[0], &outcomesA[1]

		for j := i + 1; j < len(match.Bookmakers); j++ {
			bookmakerB := &match.Bookmakers[j]
			outcomesB := bookmakerB.Markets[0].Outcomes

			// prune non-two-way bookie
			if len(outcomesB) != 2 {
				continue
			}

			oB1, oB2 := &outcomesB[0], &outcomesB[1]

			if !(oA1.Name == oB1.Name && oA2.Name == oB2.Name) {
				fmt.Println("Bookmaker names don't match into a pair")
				continue
			}

			checkArb(match, []*Outcome{oA1, oB2}, arbs)
			checkArb(match, []*Outcome{oA2, oB1}, arbs)
		}
	}
}

func findThreeWayMatchArbs(match *Match, arbs *[]Arb) {
	for i := 0; i < len(match.Bookmakers); i++ {
		bookmakerA := &match.Bookmakers[i]
		outcomesA := bookmakerA.Markets[0].Outcomes

		// prune non-two-way bookie
		if len(outcomesA) != 3 {
			continue
		}

		oA1, oA2, oA3 := &outcomesA[0], &outcomesA[1], &outcomesA[2]

		for j := i + 1; j < len(match.Bookmakers); j++ {
			bookmakerB := &match.Bookmakers[j]
			outcomesB := bookmakerB.Markets[0].Outcomes

			// prune non-two-way bookie
			if len(outcomesB) != 3 {
				continue
			}

			oB1, oB2, oB3 := &outcomesB[0], &outcomesB[1], &outcomesB[2]

			for k := j + 1; k < len(match.Bookmakers); k++ {
				bookmakerC := &match.Bookmakers[k]
				outcomesC := bookmakerC.Markets[0].Outcomes

				if len(outcomesC) != 3 {
					continue
				}

				oC1, oC2, oC3 := &outcomesC[0], &outcomesC[1], &outcomesC[2]

				if !((oA1.Name == oB1.Name && oB1.Name == oC1.Name) &&
					(oA2.Name == oB2.Name && oB2.Name == oC2.Name) &&
					(oA3.Name == oB3.Name && oB3.Name == oC3.Name)) {
					println("Bookmaker names don't match into triple")
					continue
				}

				checkArb(match, []*Outcome{oA1, oB2, oC3}, arbs)
				checkArb(match, []*Outcome{oA1, oB3, oC2}, arbs)
				checkArb(match, []*Outcome{oA2, oB1, oC3}, arbs)
				checkArb(match, []*Outcome{oA2, oB3, oC1}, arbs)
				checkArb(match, []*Outcome{oA3, oB1, oC2}, arbs)
				checkArb(match, []*Outcome{oA3, oB2, oC1}, arbs)
			}
		}
	}
}
