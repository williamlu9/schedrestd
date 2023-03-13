package cron

import (
	"schedrestd/common"
	"schedrestd/common/kvdb"
	"schedrestd/common/logger"
	"github.com/robfig/cron"
	"go.uber.org/fx"
	"strconv"
	"time"
)

// NewCron ...
func NewCron(db *kvdb.KVStore) *cron.Cron {
	c := cron.New()

	// Remove expired token form kv store
	// Run this task 4am every day
	c.AddFunc("0 0 4 * * *", func() {
		// At least, keep 24 hours for expired token
		value := time.Now().Unix() - 3600 * 24
		err := db.RemoveSmallValues(common.BoltDBJWTTable, strconv.FormatInt(value, 10))
		if err != nil {
			logger.GetDefault().Errorf("Failed to remove expired token from db: %v", err.Error())
		} else {
			logger.GetDefault().Info("Successfully removed expired token from db.")
		}
	})

	return c
}

var Module = fx.Options(fx.Provide(NewCron))
