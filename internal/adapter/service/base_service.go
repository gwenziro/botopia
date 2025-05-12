package service

import (
	"context"
	"time"

	"github.com/gwenziro/botopia/internal/domain/errors"
	"github.com/gwenziro/botopia/internal/infrastructure/logger"
)

// BaseService adalah struktur dasar untuk semua service
type BaseService struct {
	log *logger.Logger
}

// NewBaseService membuat instance BaseService baru
func NewBaseService(moduleName string, logLevel logger.LogLevel) *BaseService {
	return &BaseService{
		log: logger.New(moduleName, logLevel, true),
	}
}

// WrapWithTimeout membungkus eksekusi fungsi dengan timeout
func (s *BaseService) WrapWithTimeout(
	ctx context.Context,
	duration time.Duration,
	f func(ctx context.Context) (interface{}, error),
) (interface{}, error) {
	// Buat context dengan timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	// Buat channel untuk hasil dan error
	resultCh := make(chan interface{}, 1)
	errCh := make(chan error, 1)

	// Jalankan fungsi dalam goroutine
	go func() {
		result, err := f(timeoutCtx)
		if err != nil {
			errCh <- err
			return
		}
		resultCh <- result
	}()

	// Tunggu hasil atau timeout
	select {
	case <-timeoutCtx.Done():
		// Timeout atau dibatalkan
		if timeoutCtx.Err() == context.DeadlineExceeded {
			return nil, errors.NewTimeoutError("operasi melebihi batas waktu")
		}
		return nil, timeoutCtx.Err()
	case err := <-errCh:
		// Terjadi error
		return nil, err
	case result := <-resultCh:
		// Berhasil
		return result, nil
	}
}

// GetLogger mengembalikan logger
func (s *BaseService) GetLogger() *logger.Logger {
	return s.log
}
