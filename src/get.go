package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"go.mongodb.org/mongo-driver/bson"
)

func getSports() []Sport {
	var sports []Sport

	// find source type - this will be different later
	var source int
	sportCnt, _ := db.sports.CountDocuments(context.TODO(), bson.M{})
	if sportCnt <= 0 {
		source = SOURCE_API
	} else {
		source = SOURCE_DB
	}

	switch source {
	case SOURCE_DB:
		sports = getSportsFromDb()
		break
	case SOURCE_API:
		sports = getSportsFromApi()
		db.writeSports(sports)
		break
	case SOURCE_FILE:
		filePath := config.OutputDir + "/sports.json"
		sports = getSportsFromFile(filePath)
		break
	default:
		fmt.Println("Unknown source type:", source)
		return nil
	}

	return sports
}

func getSportsFromDb() []Sport {
	fmt.Println("Sports from db")
	sports, err := db.readSports()
	if err != nil {
		fmt.Println("Error reading sports from db")
		return nil
	}
	return sports
}

func getSportsFromApi() []Sport {
	fmt.Println("Sports from api")
	getSportsUrl := buildGetSportsUrl()

	resp, err := http.Get(getSportsUrl)
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var sports []Sport
	err = json.Unmarshal(body, &sports)
	if err != nil {
		panic(err)
	}

	return sports
}

func getSportsFromFile(filePath string) []Sport {
	if !fileExists(filePath) {
		fmt.Println("Sports file doesn't exist:", filePath)
		return nil
	}

	var sports []Sport

	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", filePath)
		return nil
	}

	err = json.Unmarshal(data, &sports)
	if err != nil {
		fmt.Println("Error decoding file as json:", err)
		return nil
	}

	return sports
}

// func writeSportsToFile(sportsJsonRawBytes []byte, outputFilePath string) {
// 	// make json pretty
// 	var prettyBody bytes.Buffer
// 	err := json.Indent(&prettyBody, sportsJsonRawBytes, "", "	")
// 	if err != nil {
// 		panic(err)
// 	}

// 	// write out to file
// 	err = os.WriteFile(outputFilePath, prettyBody.Bytes(), 0644)
// 	if err != nil {
// 		panic(err)
// 	}
// }

func buildGetSportsUrl() string {
	baseURL := config.BaseURL

	u, err := url.Parse(baseURL)
	if err != nil {
		panic(err)
	}

	apiKey := config.findApiKey()

	// add parameters
	q := u.Query()
	q.Set("apiKey", apiKey)
	u.RawQuery = q.Encode()

	return u.String()
}

func getSportMatches(sportKey string) []Match {
	var matches []Match

	// check if match list for `sportKey` exists
	filter := bson.M{"sportkey": sportKey}
	count, _ := db.matches.CountDocuments(context.TODO(), filter)
	sportKeyExists := count > 0

	var source int
	if sportKeyExists {
		source = SOURCE_DB
	} else {
		source = SOURCE_API
	}

	switch source {
	case SOURCE_DB:
		matches = getSportMatchesFromDb(sportKey)
		break
	case SOURCE_API:
		matches = getSportMatchesFromApi(sportKey)
		postProcessMatches(matches)
		db.writeSportMatches(sportKey, matches)
		break
	case SOURCE_FILE:
		filePath := config.OutputDir + "/" + config.DefaultParams.Regions + "_" + sportKey + "_odds.json"
		matches = getSportMatchesFromFile(filePath)
		break
	default:
		fmt.Println("Unknown source type:", source)
	}

	return matches
}

func getSportMatchesFromDb(sportKey string) []Match {
	fmt.Println("Matches from DB:", sportKey)
	matches, err := db.readSportMatches(sportKey)
	if err != nil {
		fmt.Println("Error reading sport odds from db:", sportKey)
	}
	return matches
}

func getSportMatchesFromApi(sportKey string) []Match {
	fmt.Println("Matches from API:", sportKey)
	getSportMatchesUrl := buildGetSportMatchesUrl(sportKey)

	resp, err := http.Get(getSportMatchesUrl)
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// TEST
	// writeSportOddsToFile(body, "out/test.json")

	var matches []Match
	err = json.Unmarshal(body, &matches)
	if err != nil {
		panic(err)
	}

	return matches
}

func getSportMatchesFromFile(filePath string) []Match {
	var matches []Match

	if !fileExists(filePath) {
		fmt.Println("Sport odds file doesn't exist:", filePath)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file: " + filePath)
		return nil
	}

	err = json.Unmarshal(data, &matches)
	if err != nil {
		fmt.Println("Error decoding file as json: ", err)
		return nil
	}

	return matches
}

func writeSportOddsToFile(sportOddsJsonRawBytes []byte, outputFilePath string) {
	// make json pretty
	var prettyBody bytes.Buffer
	err := json.Indent(&prettyBody, sportOddsJsonRawBytes, "", "	")
	if err != nil {
		panic(err)
	}

	// write out to file
	err = os.WriteFile(outputFilePath, prettyBody.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}

func buildGetSportMatchesUrl(sportKey string) string {
	// build base url + endpoint
	baseURL := config.BaseURL
	endpoint := fmt.Sprintf("%s/odds", sportKey)

	u, err := url.Parse(fmt.Sprintf("%s/%s", baseURL, endpoint))
	if err != nil {
		panic(err)
	}

	apiKey := config.findApiKey()
	regions := config.DefaultParams.Regions
	markets := config.DefaultParams.Markets
	oddsFormat := config.DefaultParams.OddsFormat

	// add parameters
	q := u.Query()
	q.Set("apiKey", apiKey)
	q.Set("regions", regions)
	q.Set("markets", markets)
	q.Set("oddsFormat", oddsFormat)
	u.RawQuery = q.Encode()

	return u.String()
}

// Add in `BookmakerKey` to each Outcome object
func postProcessMatches(matches []Match) {
	for mi := range matches {
		match := &matches[mi]
		for bi := range match.Bookmakers {
			bookmaker := &match.Bookmakers[bi]
			if len(bookmaker.Markets) == 0 {
				continue
			}
			for oi := range bookmaker.Markets[0].Outcomes {
				outcome := &bookmaker.Markets[0].Outcomes[oi]

				outcome.BookmakerKey = bookmaker.BookmakerKey
				outcome.BookmakerTitle = bookmaker.BookmakerTitle
				outcome.BookmakerUrl = config.findBookmakerUrl(bookmaker.BookmakerKey)
			}
		}
	}
}
