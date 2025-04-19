//go:build windows
// +build windows

// Package service is helper for run process as service in windows
package service

import (
	"github.com/rs/zerolog/log"

	"golang.org/x/sys/windows/svc"
)

// SCM command handler (start, stop)
type windowsServiceHandler struct {
	runApp func()
}

func (h *windowsServiceHandler) Execute(_ []string, r <-chan svc.ChangeRequest, s chan<- svc.Status) (bool, uint32) {
	s <- svc.Status{State: svc.StartPending}
	go h.runApp()
	s <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop}

	// catch SCM requests
	for c := range r {
		switch c.Cmd {
		case svc.Stop, svc.Shutdown:
			s <- svc.Status{State: svc.StopPending}
			return false, 0
		}
	}

	return false, 0
}

// IsServiceMode is run as windows service
func IsServiceMode() bool {
	isService, _ := svc.IsWindowsService()
	return isService
}

// RunAsService is run as windows service
func RunAsService(runApp func()) {
	err := svc.Run("dayz-exporter", &windowsServiceHandler{runApp: runApp})
	if err != nil {
		log.Fatal().Msgf("Service fail with error: %v", err)
	}
}
