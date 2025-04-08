package main

import (
	"errors"
	"fmt"
	"os"
)

type Sport struct {
	SportKey     string `json:"key"`
	HasOutrights bool   `json:"has_outrights"`
}

// Together, Match, Bookmaker, Market and Outcome
// objects define the structure of a Match returned by
// the odds API.
//
// JSON parses them directly.
type Match struct {
	MatchId      string      `json:"id"`
	SportKey     string      `json:"sport_key"`
	SportTitle   string      `json:"sport_title"`
	CommenceTime string      `json:"commence_time"`
	HomeTeam     string      `json:"home_team"`
	AwayTeam     string      `json:"away_team"`
	Bookmakers   []Bookmaker `json:"bookmakers"`
}

type Bookmaker struct {
	BookmakerKey   string   `json:"key"`
	BookmakerTitle string   `json:"title"`
	Markets        []Market `json:"markets"`
}

type Market struct {
	Outcomes []Outcome `json:"outcomes"`
}

type Outcome struct {
	BookmakerKey   string
	BookmakerTitle string
	BookmakerUrl   string
	Name           string
	Price          float64
}

// Represents an arbitrage opportunity
// i.e. sum of 1/o guaranteed to be < 1
type Arb struct {
	MatchId      string     `json:"matchid"`
	SportKey     string     `json:"sportkey"`
	SportTitle   string     `json:"sporttitle"`
	CommenceTime string     `json:"commencetime"`
	HomeTeam     string     `json:"hometeam"`
	AwayTeam     string     `json:"awayteam"`
	Outcomes     []*Outcome `json:"outcomes"`
	R            float64    `json:"r"`
}

func (arb *Arb) toString() {
	fmt.Printf("arb: %s %.5f (%.5f)", arb.SportKey, arb.R, (1-arb.R)*100)
	for _, o := range arb.Outcomes {
		fmt.Print(" ", o.BookmakerKey, " ", o.Name, " ", o.Price)
	}
	fmt.Println()
}

func findMatchArbs(match *Match, arbs *[]Arb) {
	// calculate most frequent number of outcomes
	numOutcomesFreqs := make(map[int]int)
	for i := range match.Bookmakers {
		bookmaker := &match.Bookmakers[i]
		market := &bookmaker.Markets[0] // NOTE: always take first market, rarely more than one
		numOutcomesFreqs[len(market.Outcomes)]++
	}
	mostFrequentNumOutcomes := findMaxKey(numOutcomesFreqs)

	switch mostFrequentNumOutcomes {
	case 0:
		// fmt.Println("No odds")
	case 2:
		findTwoWayMatchArbs(match, arbs)
		break
	case 3:
		findThreeWayMatchArbs(match, arbs)
		break
	default:
		fmt.Println("Only support 2 and 3-outcome matches:", mostFrequentNumOutcomes)
		return
	}
}

// global config
var config *Config

// global db
var db *Database

// sources
const SOURCE_DB int = 0
const SOURCE_API int = 1
const SOURCE_FILE int = 2

func findArbs() []Arb {
	// get sports list
	var sports []Sport

	sports = getSports()

	arbs := make([]Arb, 0)
	cnt := 0
	for _, sport := range sports {
		if sport.HasOutrights {
			continue
		}
		if cnt > 100 {
			break
		}
		cnt++

		sportKey := sport.SportKey

		// get odds for this sportkey and write to db
		matches := getSportMatches(sportKey)

		for i := range matches {
			match := &matches[i]
			findMatchArbs(match, &arbs)
		}
	}

	return arbs
}

func getArbs() []Arb {
	arbs, err := db.readArbs()
	if err != nil {
		panic(err)
	}
	return arbs
}

func showArbs(arbs []Arb) {
	for _, arb := range arbs {
		arb.toString()
	}
}

type CliArgs struct {
	runningInDocker bool
}

func usage() string {
	return "Usage: arb [--docker]"
}

func parseCliArgs() (CliArgs, error) {
	var cliArgs CliArgs
	nArgs := len(os.Args)

	if nArgs > 2 {
		return cliArgs, errors.New(usage())
	}

	switch nArgs {
	case 1:
		cliArgs.runningInDocker = false
	case 2:
		if os.Args[1] != "--docker" {
			return cliArgs, errors.New(usage())
		}
		cliArgs.runningInDocker = true
		break
	default:
		return cliArgs, errors.New(usage())
	}

	return cliArgs, nil
}

func main() {
	cliArgs, err := parseCliArgs()
	if err != nil {
		fmt.Println(err)
		return
	}

	initConfig(&cliArgs)
	initDb()

	// arbs := findArbs()
	// db.writeArbs(arbs)

	// arbs := getArbs()
	// showArbs(arbs)
	// arbs[0].toString()

	// println("Hello")

	startAPIServer()
}
