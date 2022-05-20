package main

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init_db() {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@clusterx.xxxxx.mongodb.net/Comet?retryWrites=true&w=majority", db_user, db_pass))
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err.Error())
	}

	client.Connect(ctx)
	collection = client.Database("Comet").Collection("Users")
}

func load_users() {
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		panic(err.Error())
	}

	for cursor.Next(ctx) {
		u := user{}

		err := cursor.Decode(&u)
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println(u)

		for _, t := range themes {
			if t.Name == u.ThemeName {
				u.Theme = t
			}
		}

		users = append(users, u)
		fmt.Println(u.Username)
	}
}

func set_ban(username string, banned bool) {
	cursor, err := collection.Find(ctx, bson.M{"username": username})
	if err != nil {
		panic(err.Error())
	}

	for cursor.Next(ctx) {
		u := user{}

		err := cursor.Decode(&u)
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println(u)
		collection.UpdateOne(ctx, bson.M{"_id": u.ID}, bson.M{"$set": bson.M{"banned": banned}})
	}
}

func set_theme(username string, theme string) {
	cursor, err := collection.Find(ctx, bson.M{"username": username})
	if err != nil {
		panic(err.Error())
	}

	for cursor.Next(ctx) {
		u := user{}

		err := cursor.Decode(&u)
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println(u)
		collection.UpdateOne(ctx, bson.M{"_id": u.ID}, bson.M{"$set": bson.M{"tfx_theme": theme}})
	}
}
