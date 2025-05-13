package main

import (
	"context"
	"log"

	"github.com/ssssunat/hotel-reservation/db"
	"github.com/ssssunat/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	ctx        = context.Background()
)

func seedHotel(name string, location string, rating int) error {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating: rating,
	}
	rooms := []types.Room{
		{
			Size:      "small",
			Price: 99.9,
		},
		{
			Size:      "normal",
			Price: 122.9,
		},
		{
			Size:      "kingsize",
			Price: 1222.9,
		},
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		_, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func main() {
	seedHotel("Belucia", "France", 3)
	seedHotel("The cozy hotel", "Netherladns", 4)
	seedHotel("Dont die in your sleep", "London", 1)
}

func init() {
	var err error
	ctx := context.Background()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
}
