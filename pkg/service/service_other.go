//go:build !windows
// +build !windows

package service

import log "github.com/sirupsen/logrus"

// always return false on all platforms except windows
func IsServiceMode() bool {
	return false
}

// just fail on all platforms except windows
func RunAsService(runApp func()) {
	log.Fatalf("Services not supported on this platform")
}
