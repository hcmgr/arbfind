package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	// db uri
	uri string

	// base mongodb object
	dbConn *mongo.Database

	// collections
	sports  *mongo.Collection
	matches *mongo.Collection
	arbs    *mongo.Collection
}

func initDb() {
	db = &Database{uri: config.MongoDbUri}
	err := db.connect()
	if err != nil {
		panic(err)
	}
}

func (db *Database) connect() error {
	clientOpts := options.Client().ApplyURI(db.uri)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		return err
	}

	// ping to ensure connection was successful
	if err := client.Ping(context.TODO(), nil); err != nil {
		return err
	}

	db.dbConn = client.Database("arb")

	db.sports = db.dbConn.Collection("sports")
	db.matches = db.dbConn.Collection("matches")
	db.arbs = db.dbConn.Collection("arbs")

	return nil
}

func (db *Database) readSports() ([]Sport, error) {
	cursor, err := db.sports.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var sports []Sport
	for cursor.Next(context.TODO()) {
		var sport Sport
		if err := cursor.Decode(&sport); err != nil {
			return nil, err
		}
		sports = append(sports, sport)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return sports, nil
}

func (db *Database) writeSports(sports []Sport) error {
	for _, sport := range sports {
		_, err := db.sports.InsertOne(context.TODO(), sport)
		if err != nil {
			return err
		}
	}
	fmt.Println("Inserted sports successfully")
	return nil
}

func (db *Database) readSportMatches(sportKey string) ([]Match, error) {
	var doc struct {
		SportKey string  `bson:"sportkey"`
		Matches  []Match `bson:"matches"`
	}

	err := db.matches.FindOne(context.TODO(), bson.M{"sportkey": sportKey}).Decode(&doc)
	if err != nil {
		return nil, err
	}

	return doc.Matches, nil
}

func (db *Database) writeSportMatches(sportKey string, matches []Match) error {
	doc := bson.M{
		"sportkey": sportKey,
		"matches":  matches,
	}
	_, err := db.matches.InsertOne(context.TODO(), doc)
	if err != nil {
		return err
	}

	fmt.Println("Inserted matches successfully")
	return nil
}

func (db *Database) readArbs() ([]Arb, error) {
	cursor, err := db.arbs.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var arbs []Arb
	for cursor.Next(context.TODO()) {
		var arb Arb
		if err := cursor.Decode(&arb); err != nil {
			return nil, err
		}
		arbs = append(arbs, arb)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return arbs, nil
}

func (db *Database) writeArbs(arbs []Arb) error {
	for _, arb := range arbs {
		_, err := db.arbs.InsertOne(context.TODO(), arb)
		if err != nil {
			return err
		}
	}
	fmt.Println("Inserted arbs successfully")
	return nil
}
