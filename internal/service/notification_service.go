package service

import (
	"context"
	"fmt"

	"github.com/juliannGabrielDev/intelfy-api/internal/repository"
	"github.com/juliannGabrielDev/intelfy-api/internal/ws"
	"github.com/juliannGabrielDev/intelfy-api/pkg/nanoid"
)

type NotificationService struct {
	repo *repository.Queries
	hub  *ws.Hub
}

func NewNotificationService(repo *repository.Queries, hub *ws.Hub) *NotificationService {
	return &NotificationService{
		repo: repo,
		hub:  hub,
	}
}

func (s *NotificationService) NotifyFollowers(ctx context.Context, artistID, title, message string) error {
	followers, err := s.repo.GetFollowersByArtistID(ctx, artistID)
	if err != nil {
		return err
	}

	for _, followerID := range followers {
		id, _ := nanoid.GenerateID()
		notification := repository.Notification{
			ID:      id,
			UserID:  followerID,
			Title:   title,
			Message: message,
		}

		err := s.repo.CreateNotification(ctx, repository.CreateNotificationParams{
			ID:      notification.ID,
			UserID:  notification.UserID,
			Title:   notification.Title,
			Message: notification.Message,
		})
		if err != nil {
			fmt.Printf("Error creating notification for user %s: %v\n", followerID, err)
			continue
		}

		// Push via WebSocket if connected
		if s.hub != nil {
			s.hub.SendToUser(followerID, notification)
		}
	}

	return nil
}
