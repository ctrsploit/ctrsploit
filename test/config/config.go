package config

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/ssst0n3/awesome_libs/log"
)

const (
	EnvInDocker = "TEST_IN_DOCKER"
)

var InDocker bool

func Init() {
	InDocker, _ = strconv.ParseBool(os.Getenv(EnvInDocker))
	log.Logger.Level = logrus.DebugLevel
}
