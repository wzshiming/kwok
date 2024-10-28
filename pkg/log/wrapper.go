/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package log

import (
	"fmt"
	"log/slog" //nolint:depguard
	"strconv"
	"strings"
)

// Level is the logging level.
type Level = slog.Level

// The following is Level definitions copied from slog.
const (
	LevelDebug Level = slog.LevelDebug
	LevelInfo  Level = slog.LevelInfo
	LevelWarn  Level = slog.LevelWarn
	LevelError Level = slog.LevelError
)

func wrapSlog(handler slog.Handler, level Level) *Logger {
	return &Logger{
		Logger: slog.New(handler),
		level:  level,
	}
}

// Logger is a wrapper around slog.Handler.
type Logger struct {
	*slog.Logger
	level Level
}

// Level returns the level of the logger
func (l *Logger) Level() Level {
	return l.level
}

// With returns a Logger that includes the given attributes in each output operation.
func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		Logger: l.Logger.With(args...),
		level:  l.level,
	}
}

// WithGroup returns a Logger that starts a group, if name is non-empty.
func (l *Logger) WithGroup(name string) *Logger {
	return &Logger{
		Logger: l.Logger.WithGroup(name),
		level:  l.level,
	}
}

// ParseLevel parses a level string.
func ParseLevel(s string) (l Level, err error) {
	name := s
	offsetStr := ""
	i := strings.IndexAny(s, "+-")
	if i > 0 {
		name = s[:i]
		offsetStr = s[i:]
	} else if i == 0 ||
		(name[0] >= '0' && name[0] <= '9') {
		name = "INFO"
		offsetStr = s
	}

	switch strings.ToUpper(name) {
	case "DEBUG":
		l = LevelDebug
	case "INFO":
		l = LevelInfo
	case "WARN":
		l = LevelWarn
	case "ERROR":
		l = LevelError
	default:
		return 0, fmt.Errorf("ParseLevel %q: invalid level name", s)
	}

	if offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return 0, fmt.Errorf("ParseLevel %q: invalid offset: %w", s, err)
		}
		l += Level(offset)
	}

	return l, nil
}
