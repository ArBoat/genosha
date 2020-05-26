package dao

import (
	"context"
	"fmt"
	"genosha/db"
	"log"
	//"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/mongo/options"
	//"log"
)

var mdb = db.MongoDB

func init() {
	//testMongoDB()
}

type Trainer struct {
	Name string
	Age  int
	City string
}

func testMongoDB() {
	collection := mdb.Database("test").Collection("trainers")
	ash := Trainer{"Ash", 10, "Pallet Town"}
	misty := Trainer{"Misty", 10, "Cerulean City"}
	brock := Trainer{"Brock", 15, "Pewter City"}
	// insert
	insertResult, err := collection.InsertOne(context.TODO(), ash)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	trainers := []interface{}{misty, brock}

	insertManyResult, err := collection.InsertMany(context.TODO(), trainers)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)

	////update
	//filter := bson.D{{"name", "Ash"}}
	//update := bson.D{
	//	{"$inc", bson.D{
	//		{"age", 1},
	//	}},
	//}
	//updateResult, err := collection.UpdateMany(context.TODO(), filter, update)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	//var result Trainer
	//err := collection.FindOne(context.TODO(), filter).Decode(&result)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(result)
	//fmt.Printf("Found a single document: %+v\n", result)
	//
	//findOptions := options.Find()
	//findOptions.SetLimit(2)
	//var results []Trainer
	//cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//for cur.Next(context.TODO()) {
	//
	//	// create a value into which the single document can be decoded
	//	var elem Trainer
	//	err := cur.Decode(&elem)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	results = append(results, elem)
	//}
	//if err := cur.Err(); err != nil {
	//	log.Fatal(err)
	//}
	//
	//// Close the cursor once finished
	//cur.Close(context.TODO())
	//fmt.Println(results)
	//fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)
	//
	//deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{{}})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
	err = collection.Drop(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}
