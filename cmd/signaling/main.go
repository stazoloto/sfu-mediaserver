package main

import (
	"log"
	"net/http"

	"github.com/stazoloto/sfu-mediaserver/internal/sfu"
	"github.com/stazoloto/sfu-mediaserver/internal/signaling/interfaceadapters/controllers"
	"github.com/stazoloto/sfu-mediaserver/internal/signaling/interfaceadapters/presenters"
	"github.com/stazoloto/sfu-mediaserver/internal/signaling/interfaceadapters/repositories"
	"github.com/stazoloto/sfu-mediaserver/internal/signaling/usecase"
	"github.com/stazoloto/sfu-mediaserver/pkg/ws"
)

func main() {
	roomRepo := repositories.NewRoomRepo()
	hub := ws.NewHub()
	presenter := presenters.NewWSPresenter(hub)
	sfu := sfu.NewSFU()

	interactor := usecase.NewInteractor(roomRepo, presenter, sfu)
	controller := controllers.NewWSController(interactor)
	hub.SetController(controller)

	hub.SetOnDisconnect(interactor.Disconnect)
	http.Handle("/ws", hub)

	log.Println("server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
