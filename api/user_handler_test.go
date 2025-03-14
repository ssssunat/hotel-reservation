package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/ssssunat/hotel-reservation/db"
	"github.com/ssssunat/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbname = "hotel-reservation-test"
)

type testdb struct {
	db.UserStore
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatal(err)
	}
	testdburi := os.Getenv("MONGO_DB_URL_TEST")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testdburi))
	if err != nil {
		log.Fatal(err)
	}

	return &testdb{
		UserStore: db.NewMongoUserStore(client),
	}
}

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userhandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userhandler.HandlePostUser)

	params := types.CreateUserParams{
		Email:     "some@foo.com",
		FirstName: "James",
		LastName:  "Foo",
		Password:  "asdasdasqweqw",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Errorf(err.Error())
	}

	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)
	if len(user.ID) == 0 {
		t.Errorf("expecting a user id to be set")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("expecting encrypted password not included into json respose")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected firstname %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected firstname %s but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected email %s but got %s", params.Email, user.Email)
	}
}

// TODO GET DELETE UPDATE handlers

func TestGetUsers(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Get("/", userHandler.HandleGetUsers)

	users := []types.User{
		{
			Email:     "john@doe.com",
			FirstName: "John",
			LastName:  "Doe",
		},
		{
			Email:     "jane@doe.com",
			FirstName: "Jane",
			LastName:  "Doe",
		},
	}

	for _, user := range users {
		_, err := tdb.UserStore.InsertUser(context.Background(), &user)
		if err != nil {
			t.Fatalf("failed to create test user:%v", err)
		}
	}

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Content-Type", "application/json")

	resp, _ := app.Test(req)
	fmt.Println(resp.Status)

	var getUsers []types.User
	if err := json.NewDecoder(resp.Body).Decode(&getUsers); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	for i, expected := range users {
		user := getUsers[i]
		if user.Email != expected.Email {
			t.Errorf("Expected email %s, but got %s", expected.Email, user.Email)
		}
		if user.FirstName != expected.FirstName {
			t.Errorf("Expected email %s, but got %s", expected.FirstName, user.FirstName)
		}
	}
	fmt.Println(getUsers)
}

