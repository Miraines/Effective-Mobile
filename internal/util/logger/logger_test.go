package logger_test

import (
	"Effective-Mobile/internal/util/logger"
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestNewLevels(t *testing.T) {
	t.Parallel()
	cases := []struct {
		lvl    string
		expect zapcore.Level
	}{{"debug", zapcore.DebugLevel}, {"info", zapcore.InfoLevel}, {"warn", zapcore.WarnLevel}, {"error", zapcore.ErrorLevel}, {"fatal", zapcore.FatalLevel}}

	for _, c := range cases {
		l, err := logger.New(c.lvl)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !l.Core().Enabled(c.expect) {
			t.Errorf("level %s should enable %s", c.lvl, c.expect)
		}
	}
}
