package namespace

import (
	"github.com/ctrsploit/ctrsploit/log"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const (
	CtrsploitTestNetNsHost = "CTRSPLOIT_TEST_NET_NS_HOST"
)

// TODO: test in e2e
func TestInoArbitrator_IsNetworkNamespaceInoBetweenProcInoList(t *testing.T) {
	log.Logger.Level = logrus.DebugLevel
	inoArbitrator, err := NewInoArbitrator()
	assert.NoError(t, err)
	netns := Namespace{
		Name:        "net",
		Path:        "/proc/self/ns/net",
		Type:        TypeNamespaceTypeNetwork,
		InodeNumber: 4026532726,
	}
	host := inoArbitrator.IsNetworkNamespaceInoBetweenProcInoList(netns)
	assert.Equal(t, os.Getenv(CtrsploitTestNetNsHost) == "true", host)
}

func TestInoArbitrator_GuessNetworkNamespaceInitialIno(t *testing.T) {
	log.Logger.Level = logrus.DebugLevel
	inoArbitrator, err := NewInoArbitrator()
	assert.NoError(t, err)
	initialIno := inoArbitrator.GuessNetworkNamespaceInitialIno()
	log.Logger.Debug(initialIno)
}

func TestInoArbitrator_IsNetworkNamespaceInoBetweenTwoAdjacentMissingIno(t *testing.T) {
	log.Logger.Level = logrus.DebugLevel
	inoArbitrator, err := NewInoArbitrator()
	assert.NoError(t, err)
	netns := Namespace{
		Name:        "net",
		Path:        "/proc/self/ns/net",
		Type:        TypeNamespaceTypeNetwork,
		InodeNumber: 4026532726,
	}
	host := inoArbitrator.IsNetworkNamespaceInoBetweenTwoAdjacentMissingIno(netns)
	assert.Equal(t, os.Getenv(CtrsploitTestNetNsHost) == "true", host)
}
