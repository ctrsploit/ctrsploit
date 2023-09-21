package util

import (
	"github.com/ctrsploit/ctrsploit/internal/log"
	"github.com/sirupsen/logrus"
	log2 "github.com/ssst0n3/awesome_libs/log"
)

func Debug() {
	log.Logger.SetReportCaller(true)
	log.Logger.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	log.Logger.Level = logrus.DebugLevel
	log2.Logger.SetReportCaller(true)
	log2.Logger.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	log2.Logger.Level = logrus.DebugLevel
	log.Logger.Debug("debug mode enabled")
}
