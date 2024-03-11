package slogexp

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

// ReplaceAttrFunc is a function that replaces an attribute.
type ReplaceAttrFunc func(groups []string, a slog.Attr) slog.Attr

// WrapReplaceAttrFunc wraps a ReplaceAttrFunc.
func WrapReplaceAttrFunc(fns ...ReplaceAttrFunc) ReplaceAttrFunc {
	return func(groups []string, a slog.Attr) slog.Attr {
		for _, fn := range fns {
			a = fn(groups, a)
		}
		return a
	}
}

// ReplaceAttrFuncs is a list of ReplaceAttrFuncs.
func ReplaceTimeAttr(timeFormat string) ReplaceAttrFunc {
	if timeFormat == "" {
		timeFormat = time.RFC3339Nano
	}

	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			a.Value = slog.StringValue(a.Value.Time().Format(timeFormat))
		}
		return a
	}
}

func ReplaceSourceAttr() ReplaceAttrFunc {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.SourceKey {
			source := a.Value.Any().(*slog.Source)
			items := strings.Split(source.File, "/")
			if len(items) > 3 {
				fmt.Println(items)
				if items[0] == "" {
					a.Value = slog.StringValue("/" + items[1] + "/.../" + items[len(items)-1] + ":" + strconv.Itoa(source.Line))
				} else {
					a.Value = slog.StringValue(items[0] + "/.../" + items[len(items)-1] + ":" + strconv.Itoa(source.Line))
				}

			}

		}
		return a
	}
}
