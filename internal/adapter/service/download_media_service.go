package service

import (
	"context"
	"time"

	"github.com/gwenziro/botopia/internal/domain/message"
	"github.com/gwenziro/botopia/internal/domain/repository"
	"github.com/gwenziro/botopia/internal/infrastructure/logger"
)

// DownloadMediaService implementasi layanan unduh media
type DownloadMediaService struct {
	connectionRepo repository.ConnectionRepository
	log            *logger.Logger
}

// NewDownloadMediaService membuat instance layanan unduh media baru
func NewDownloadMediaService(
	connectionRepo repository.ConnectionRepository,
	log *logger.Logger,
) *DownloadMediaService {
	return &DownloadMediaService{
		connectionRepo: connectionRepo,
		log:            log,
	}
}

// DownloadMedia mengunduh media dari pesan
func (s *DownloadMediaService) DownloadMedia(msg *message.Message) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	s.log.Info("Mengunduh media dari pesan ID: %s", msg.ID)

	// Delegasikan ke connection repository
	return s.connectionRepo.DownloadMedia(ctx, msg)
}
