package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cliOpts := options.Client().SetHosts([]string{"127.0.0.1:45000"}).SetDirect(true)
	cli, err := mongo.NewClient(cliOpts)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = cli.Connect(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cli.Disconnect(context.Background())
	trueVar := true

	cursor, err := cli.Database("test").Collection("fff").Find(context.Background(), bson.M{}, &options.FindOptions{
		ShowRecordID: &trueVar,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	for cursor.Next(context.Background()) {
		var doc bson.M
		err = cursor.Decode(&doc)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(doc)
	}
}
