package api

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ssssunat/hotel-reservation/db"
)

type BookRoomParams struct {
	FromDate time.Time `json:"fromDate`
	TillDate time.Time `json:"tillDate`
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	return nil
}
