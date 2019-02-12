package main

import (
	"strings"

	"github.com/EGT-Ukraine/gitlab-trigger-watcher/pipeline"
	"github.com/pkg/errors"
)

func schemaConverter(s string) (pipeline.Schema, error) {
	switch strings.ToLower(s) {
	case "https":
		return pipeline.HTTPS, nil
	case "http":
		return pipeline.HTTP, nil
	}

	return 0, errors.New(`wrong schema passed! Should be one of: "http" or "https"`)
}
