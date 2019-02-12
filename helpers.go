package main

import (
	"strings"

	"github.com/EGT-Ukraine/gitlab-trigger-watcher/trigger"
	"github.com/pkg/errors"
)

func schemaConverter(s string) (trigger.Schema, error) {
	switch strings.ToLower(s) {
	case "https":
		return trigger.HTTPS, nil
	case "http":
		return trigger.HTTP, nil
	}

	return 0, errors.New(`wrong schema passed! Should be one of: "http" or "https"`)
}
