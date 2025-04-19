//go:build !windows
// +build !windows

// Package service is helper for run process as service in other OS non Windows
package service

import "github.com/rs/zerolog/log"

// IsServiceMode always return false on all platforms except windows
func IsServiceMode() bool {
	return false
}

// RunAsService just fail on all platforms except windows
func RunAsService(_ func()) {
	log.Fatal().Msgf("Services not supported on this platform")
}
