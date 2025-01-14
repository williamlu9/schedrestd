package service

import (
	"schedrestd/common"
	"schedrestd/common/kvdb"
	"schedrestd/common/logger"
	"context"
	"github.com/robfig/cron"
	"go.uber.org/fx"
)

// Service ...
type Service struct {
	Cron   *cron.Cron
	Db     *kvdb.KVStore
	Logger logger.AipLogger
}

// NewService ...
func NewService(c *cron.Cron, db *kvdb.KVStore) Service {
	return Service {
		Cron:   c,
		Db:     db,
		Logger: logger.GetDefault(),
	}
}

// Initial services
func (s *Service) Initial(ctx context.Context) error {
	s.Cron.Start()
	s.Logger.Info("Cron service is started.")

	err := s.Db.Open()
	if err != nil {
		s.Logger.Error(err.Error())
		return err
	}
	s.Logger.Info("KV db is opened.")

	// Init the jwt bolt table
	err = s.Db.CreateTableIfNotExist(common.BoltDBJWTTable)
	if err != nil {
		s.Logger.Error(err.Error())
		return err
	}

	go s.closeService(ctx)

	return nil
}

func (s *Service) closeService(ctx context.Context) {
	select {
	case <- ctx.Done():
		s.Db.Close()
		s.Cron.Stop()
	}
}

// Module ...
var Module = fx.Options(fx.Provide(NewService))

