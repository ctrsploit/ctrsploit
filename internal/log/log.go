package log

import (
	"github.com/sirupsen/logrus"
	"github.com/ssst0n3/awesome_libs/log/logger"
	"os"
)

var Logger *logrus.Logger

func init() {
	Logger = logger.InitLogger("ctrsploit", os.Stdout)
	Logger.SetReportCaller(false)
}
