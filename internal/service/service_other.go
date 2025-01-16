//go:build !windows
// +build !windows

package service

import "github.com/rs/zerolog/log"

// always return false on all platforms except windows
func IsServiceMode() bool {
	return false
}

// just fail on all platforms except windows
func RunAsService(_ func()) {
	log.Fatal().Msgf("Services not supported on this platform")
}
