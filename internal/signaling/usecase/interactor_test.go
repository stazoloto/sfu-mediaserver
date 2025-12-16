package usecase

import (
	"testing"

	"github.com/stazoloto/sfu-mediaserver/internal/signaling/entities"
	"github.com/stazoloto/sfu-mediaserver/internal/signaling/interfaceadapters/repositories"
)

func TestJoinRoomAndPeers(t *testing.T) {
	roomRepo := repositories.NewRoomRepo()
	gateway := NewMockClientGateway()

	uc := NewInteractor(roomRepo, gateway)

	// Alice join
	err := uc.Handle(entities.Message{
		Type:     entities.TypeJoin,
		Room:     "room1",
		ClientID: "alice",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Проверка: alice получил peers
	toAliceGatewayMsgCount := len(gateway.Sent["alice"])

	if toAliceGatewayMsgCount != 1 {
		t.Fatalf("expected 1 message to alice, got %d", toAliceGatewayMsgCount)
	}

	msg := gateway.Sent["alice"][0]
	if msg.Type != entities.TypePeers {
		t.Fatalf("expected peers message, got %s", msg.Type)
	}

	// bob join
	err = uc.Handle(entities.Message{
		Type:     entities.TypeJoin,
		Room:     "room1",
		ClientID: "bob",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Проверка: bob получил peers
	toBobGatewayMsgCount := len(gateway.Sent["bob"])

	if toBobGatewayMsgCount != 1 {
		t.Fatalf("expected 1 message to bob, got %d", toBobGatewayMsgCount)
	}

	msg = gateway.Sent["bob"][0]
	if msg.Type != entities.TypePeers {
		t.Fatalf("expected peers message, got %s", msg.Type)
	}

	// Проверка: alice получил обновленный peers
	if len(gateway.Sent["alice"]) < 2 {
		t.Fatal("alice did not receive updated peers")
	}
}

func TestJoinWithoutRoom(t *testing.T) {
	repo := repositories.NewRoomRepo()
	gateway := NewMockClientGateway()

	uc := NewInteractor(repo, gateway)

	err := uc.Handle(entities.Message{
		Type:     entities.TypeJoin,
		ClientID: "alice",
	})

	if err != ErrMissingRoom {
		t.Fatalf("expected ErrMissingRoom, got %v", err)
	}
}
