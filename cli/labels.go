package main

import (
	"github.com/woozymasta/a2s/pkg/a2s"
	"github.com/woozymasta/dayz-exporter/pkg/bemetrics"
)

// return base labels from A2S INFO and additional extra labels
func makeLabels(info *a2s.Info, extraLabels map[string]string) bemetrics.Labels {
	labels := []bemetrics.Label{
		{Key: "server", Value: info.Name},
		{Key: "map", Value: info.Map},
		{Key: "game", Value: info.Folder},
		{Key: "os", Value: info.Environment.String()},
		{Key: "version", Value: info.Version},
	}

	for k, v := range extraLabels {
		labels = append(labels, bemetrics.Label{Key: k, Value: v})
	}

	return labels
}
