package util

import (
	"os/user"
	"strconv"
	"syscall"
	"testing"

	"golang.org/x/sys/unix"
)

func TestWithUid(t *testing.T) {
	// This test needs to be run as root
	if syscall.Geteuid() != 0 {
		t.Skip("skipping test that requires root privileges")
	}

	// Find a non-root user to switch to
	u, err := user.Lookup("nobody")
	if err != nil {
		t.Fatalf("failed to find user 'nobody': %v", err)
	}
	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		t.Fatalf("failed to parse uid: %v", err)
	}

	err = WithUid(uid, func() error {
		if syscall.Geteuid() != uid {
			t.Errorf("expected euid %d, got %d", uid, syscall.Geteuid())
		}
		return nil
	})
	if err != nil {
		t.Fatalf("WithUid failed: %v", err)
	}
	if syscall.Geteuid() != 0 {
		t.Errorf("failed to restore to original euid")
	}
}

func TestWithUidAndCaps(t *testing.T) {
	// This test needs to be run as root
	if syscall.Getuid() != 0 {
		t.Skip("skipping test that requires root privileges")
	}

	// Find a non-root user to switch to
	u, err := user.Lookup("nobody")
	if err != nil {
		t.Fatalf("failed to find user 'nobody': %v", err)
	}
	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		t.Fatalf("failed to parse uid: %v", err)
	}

	// Get original capabilities
	hdr := unix.CapUserHeader{Version: unix.LINUX_CAPABILITY_VERSION_3}
	origCaps := [2]unix.CapUserData{}
	if err := unix.Capget(&hdr, &origCaps[0]); err != nil {
		t.Fatalf("failed to get original capabilities: %v", err)
	}
	t.Logf("Original effective capabilities: %x, %x", origCaps[0].Effective, origCaps[1].Effective)

	err = WithUidAndCaps(uid, func() error {
		if syscall.Geteuid() != uid {
			t.Errorf("expected euid %d, got %d", uid, syscall.Geteuid())
		}

		// Check capabilities
		newCaps := [2]unix.CapUserData{}
		if err := unix.Capget(&hdr, &newCaps[0]); err != nil {
			t.Errorf("failed to get capabilities inside WithUidAndCaps: %v", err)
		}
		t.Logf("New effective capabilities: %x, %x", newCaps[0].Effective, newCaps[1].Effective)
		if origCaps[0].Effective != newCaps[0].Effective || origCaps[1].Effective != newCaps[1].Effective {
			t.Errorf("effective capabilities do not match. expected %v, got %v", origCaps, newCaps)
		}

		return nil
	})

	if err != nil {
		t.Fatalf("WithUidAndCaps failed: %v", err)
	}
}
