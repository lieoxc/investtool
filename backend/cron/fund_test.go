package cron

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func _TestSyncFund(t *testing.T) {
	logrus.SetLevel(logrus.WarnLevel)
	viper.SetDefault("app.chan_size", 500)
	SyncFund()
}

func _TestSyncFundManagers(t *testing.T) {
	logrus.SetLevel(logrus.WarnLevel)
	viper.SetDefault("app.chan_size", 500)
	SyncFundManagers()
}
