package main

import (
	"github.com/rumblefrog/go-a2s"
	"github.com/woozymasta/dayz-exporter/pkg/bemetrics"
)

// return base labels from A2S INFO and additional extra labels
func makeLabels(info *a2s.ServerInfo, extraLabels map[string]string) bemetrics.Labels {
	game := info.Game
	if game == "" {
		game = info.Folder
	}

	labels := []bemetrics.Label{
		{Key: "server", Value: info.Name},
		{Key: "map", Value: info.Map},
		{Key: "game", Value: game},
		{Key: "os", Value: info.ServerOS.String()},
		{Key: "version", Value: info.Version},
	}

	for k, v := range extraLabels {
		labels = append(labels, bemetrics.Label{Key: k, Value: v})
	}

	return labels
}
