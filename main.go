package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Podcast struct {
	Id     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name   string             `bson:"name" json:"name"`
	Author string             `bson:"author" json:"author"`
	Tags   []string           `bson:"tags" json:"tags"`
}

type Episode struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Author      string             `bson:"author" json:"author"`
	Podcast     Podcast            `bson:"podcast json:"podcast`
	Duration    int32              `bson:"duration" json:"duration"`
	Description string             `bson:"description" json:"description"`
}

func main() {

	uri := "mongodb://localhost:27017"
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	//closing connection automatically
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged.")

	// get all databses
	databses, err := client.ListDatabaseNames(context.TODO(), bson.M{})
	if err != nil {
		panic(err)
	}
	fmt.Println(databses)

	// create a database
	quickstartDatabase := client.Database("quickstart")

	//create a collection
	podcastsCollection := quickstartDatabase.Collection("podcasts")
	episodesCollection := quickstartDatabase.Collection("episodes")

	// // insert in a collection

	podcastData := Podcast{
		Name:   "ankit 123",
		Author: "Ankit Kumar",
		Tags:   []string{"mongodb", "nosql"},
	}

	podcatsResult, err := podcastsCollection.InsertOne(context.TODO(), podcastData)
	if err != nil {
		panic(err)
	}

	fmt.Println(podcatsResult.InsertedID)

	// insert many in a collection

	// episodesResult, err := episodesCollection.InsertMany(context.TODO(), []interface{}{
	// 	bson.D{
	// 		{"podcast", podcatsResult.InsertedID},
	// 		{"title", "episode #1"},
	// 		{"description", "hiii"},
	// 		{"duration", 25},
	// 	},
	// 	bson.D{
	// 		{"podcast", podcatsResult.InsertedID},
	// 		{"title", "episode #2"},
	// 		{"description", "hiii"},
	// 		{"duration", 32},
	// 	},
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(episodesResult.InsertedIDs)

	//Retrive Data from database
	cursor, err := episodesCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		panic(err)
	}

	// var episodes []bson.M
	// err = cursor.All(context.TODO(), &episodes)
	// if err != nil {
	// 	panic(err)
	// }

	// for _, episodes := range episodes {
	// 	fmt.Println(episodes["title"])

	// }

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var episodes bson.M
		if err = cursor.Decode(&episodes); err != nil {
			panic(err)
		}
		fmt.Println(episodes)

	}

	var podcast bson.M
	if err = podcastsCollection.FindOne(context.TODO(), bson.M{}).Decode(&podcast); err != nil {
		panic(err)
	}
	fmt.Println(podcast)

	filtercursor, err := episodesCollection.Find(context.TODO(), bson.M{"duration": 25})
	if err != nil {
		panic(err)
	}

	var episodeFiltered []bson.M
	if err = filtercursor.All(context.TODO(), &episodeFiltered); err != nil {
		panic(err)
	}
	fmt.Println(episodeFiltered)

	opts := options.Find()
	opts.SetSort((bson.D{{"duration", 1}}))
	sortCursor, err := episodesCollection.Find(context.TODO(), bson.D{
		{"duration", bson.D{
			{"$gt", 26},
		}},
	}, opts)
	if err != nil {
		panic(err)
	}

	var episodeSorted []bson.M
	if err = sortCursor.All(context.TODO(), &episodeSorted); err != nil {
		panic(err)
	}
	fmt.Println("sorted episodes================", episodeSorted)

	//update Data in database

	id, _ := primitive.ObjectIDFromHex("637c6bb18f8d7c52fb9c6340")

	result, err := podcastsCollection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.D{
		{"$set", bson.M{
			"author": "ankit yadav",
		},
		}},
	)

	if err != nil {
		panic(err)
	}
	fmt.Println(result.ModifiedCount)

	// delete from database

	resultDelete, err := episodesCollection.DeleteOne(context.TODO(), bson.M{"duration": 25})
	if err != nil {
		panic(err)
	}

	fmt.Println(resultDelete.DeletedCount)

}
