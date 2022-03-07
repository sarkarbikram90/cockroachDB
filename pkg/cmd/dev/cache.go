// Copyright 2022 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

const (
	bazelRemoteTarget = "@com_github_buchgr_bazel_remote//:bazel-remote"
	cacheCleanFlag    = "clean"
	cacheDownFlag     = "down"
	cacheResetFlag    = "reset"

	cachePidFilename = ".dev-cache.pid"
	configFilename   = "config.yml"
)

func makeCacheCmd(runE func(cmd *cobra.Command, args []string) error) *cobra.Command {
	cacheCmd := &cobra.Command{
		Use:   "cache",
		Short: "Configure and manage dev cache",
		Long:  "Configure and manage dev cache.",
		Example: `dev cache
dev cache --down
dev cache --reset
dev cache --clean`,
		Args: cobra.ExactArgs(0),
		RunE: runE,
	}
	cacheCmd.Flags().Bool(cacheCleanFlag, false, "clean cache directory")
	cacheCmd.Flags().Bool(cacheDownFlag, false, "tear down local cache server")
	cacheCmd.Flags().Bool(cacheResetFlag, false, "tear down and reboot local cache server")
	return cacheCmd
}

func (d *dev) cache(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	clean := mustGetFlagBool(cmd, cacheCleanFlag)
	down := mustGetFlagBool(cmd, cacheDownFlag)
	reset := mustGetFlagBool(cmd, cacheResetFlag)

	if clean {
		return d.cleanCache(ctx)
	}
	if down {
		return d.tearDownCache(ctx)
	}
	if reset {
		// Errors here don't really mean much, we can just ignore them.
		err := d.tearDownCache(ctx)
		if err != nil {
			log.Printf("%v\n", err)
		}
	}
	bazelRcLine, err := d.setUpCache(ctx)
	if err != nil {
		return err
	}
	_, err = d.checkPresenceInBazelRc(bazelRcLine)
	return err
}

func bazelRemoteCacheDir() (string, error) {
	const bazelRemoteCacheDir = "dev-bazel-remote"
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cacheDir, bazelRemoteCacheDir), nil
}

func getCachePidFilePath() (string, error) {
	cacheDir, err := bazelRemoteCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cacheDir, cachePidFilename), nil
}

func getCachePid() (int, error) {
	pidFile, err := getCachePidFilePath()
	if err != nil {
		return 0, err
	}
	pidStr, err := os.ReadFile(pidFile)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(pidStr))
}

func (d *dev) cacheIsUp(ctx context.Context) bool {
	pid, err := getCachePid()
	if err != nil {
		return false
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = proc.Signal(syscall.Signal(0))
	return err == nil
}

// setUpCache returns a non-nil error iff setting up the cache failed, and a
// string which is a line that should be added to ~/.bazelrc.
func (d *dev) setUpCache(ctx context.Context) (string, error) {
	if d.cacheIsUp(ctx) {
		return d.getCacheBazelrcLine(ctx)
	}

	log.Printf("Configuring cache...\n")

	bazelRemoteLoc, err := d.exec.CommandContextSilent(ctx, "bazel", "run", bazelRemoteTarget, "--run_under=//build/bazelutil/whereis")
	if err != nil {
		return "", err
	}
	bazelRemoteBinary := strings.TrimSpace(string(bazelRemoteLoc))

	// write config file unless already exists
	cacheDir, err := bazelRemoteCacheDir()
	if err != nil {
		return "", err
	}
	configFile := filepath.Join(cacheDir, configFilename)
	_, err = os.Stat(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			err := d.os.MkdirAll(filepath.Join(cacheDir, "cache"))
			if err != nil {
				return "", err
			}
			err = d.os.WriteFile(configFile, fmt.Sprintf(`# File generated by dev. You can edit this file in-place.
# See https://github.com/buchgr/bazel-remote for additional information.

dir: %s
max_size: 16
host: localhost
port: 9867
`, filepath.Join(cacheDir, "cache")))
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}
	log.Printf("Using cache configuration file at %s\n", configFile)

	cmd := exec.Command(bazelRemoteBinary, "--config_file", configFile)
	stdout, err := os.Create(filepath.Join(cacheDir, "stdout.log"))
	if err != nil {
		return "", err
	}
	cmd.Stdout = stdout
	stderr, err := os.Create(filepath.Join(cacheDir, "stderr.log"))
	if err != nil {
		return "", err
	}
	cmd.Stderr = stderr
	err = cmd.Start()
	if err != nil {
		return "", err
	}
	pid := cmd.Process.Pid
	err = cmd.Process.Release()
	if err != nil {
		return "", err
	}

	err = d.os.WriteFile(filepath.Join(cacheDir, cachePidFilename), strconv.Itoa(pid))
	if err != nil {
		return "", err
	}
	return d.getCacheBazelrcLine(ctx)
}

func (d *dev) tearDownCache(ctx context.Context) error {
	pid, err := getCachePid()
	if err != nil {
		return err
	}
	cachePidFile, err := getCachePidFilePath()
	if err != nil {
		return err
	}
	err = d.os.Remove(cachePidFile)
	if err != nil {
		return err
	}
	proc, err := os.FindProcess(pid)
	if err == nil {
		return proc.Signal(syscall.SIGTERM)
	}
	return nil
}

func (d *dev) cleanCache(ctx context.Context) error {
	if d.cacheIsUp(ctx) {
		return fmt.Errorf("cache is currently running; please run `dev cache --down`")
	}
	dir, err := bazelRemoteCacheDir()
	if err != nil {
		return err
	}
	return os.RemoveAll(filepath.Join(dir, "cache"))
}

func (d *dev) getCacheBazelrcLine(ctx context.Context) (string, error) {
	cacheDir, err := bazelRemoteCacheDir()
	if err != nil {
		return "", err
	}
	configFile := filepath.Join(cacheDir, configFilename)
	// We "should" be using a YAML parser for this, but who's going to stop me?
	configFileContents, err := d.os.ReadFile(configFile)
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(configFileContents, "\n") {
		if strings.HasPrefix(line, "port:") {
			port := strings.TrimSpace(strings.Split(line, ":")[1])
			return fmt.Sprintf("build --remote_cache=http://127.0.0.1:%s", port), nil
		}
	}
	return "", fmt.Errorf("could not determine what to add to ~/.bazelrc to enable cache")
}
