package di

import (
	"github.com/gwenziro/botopia/internal/adapter/service"
	"github.com/gwenziro/botopia/internal/domain/message"
)

// initMediaServices menginisialisasi layanan media
func (c *Container) initMediaServices() {
	// Download media service
	downloadMediaService := service.NewDownloadMediaService(
		c.connectionRepository,
		c.log,
	)

	// Register service di domain layer
	message.SetDownloadMediaService(downloadMediaService)

	c.log.Info("Media services berhasil diinisialisasi")
}
