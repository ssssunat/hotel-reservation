package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/ssssunat/hotel-reservation/api"
	"github.com/ssssunat/hotel-reservation/api/middleware"
	"github.com/ssssunat/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "The listen addr of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	// handler init
	var (
		hotelStore = db.NewMongoHotelStore(client)
		roomStore = db.NewMongoRoomStore(client, hotelStore)
		userStore = db.NewMongoUserStore(client)
		store = &db.Store{
			Hotel: hotelStore,
			Room: roomStore,
			User: userStore,
		}
		userHandler = api.NewUserHandler(userStore)
		hotelHandler = api.NewHotelHandler(store)
		app = fiber.New(config)
		apiv1 = app.Group("api/v1", middleware.JWTAuthentication)
	
	)

	// user handlers
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)

	//hotel handlers
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)
	app.Listen(*listenAddr)
}

