package config

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

const (
	DefaultLogLevel      = "info"
	DefaultPort          = 9090
	DefaultBindInterface = "0.0.0.0"
	AppName              = "rpi_export"
)

type Config struct {
	AtomLogger       zap.AtomicLevel
	LogLevel         string
	BindInterface    string
	Port             int
	HostNameOverride string
	MailboxDebug     bool

	Address string
}

func InitLogger(c *Config) {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	encoder := zapcore.NewJSONEncoder(config)
	atom := zap.NewAtomicLevel()
	c.AtomLogger = atom
	logr := zap.New(zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), atom))
	zap.ReplaceGlobals(logr)
}

func setLogLevel(c *Config) error {
	var level zapcore.Level
	err := (&level).UnmarshalText([]byte(c.LogLevel))

	if err != nil {
		return err
	}

	c.AtomLogger.SetLevel(level)

	return nil
}

func ValidateConfig(c *Config) error {
	if err := setLogLevel(c); err != nil {
		return err
	}

	c.Address = fmt.Sprintf("%s:%d", c.BindInterface, c.Port)

	return nil
}
