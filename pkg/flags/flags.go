package flags

import (
	"flag"
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	logLevel    = "log-level"
	logEncoding = "log-encoding" // json or console
	configFile  = "config-file"
)

// AddConfigFileFlag adds flags for ExternalConfFlags
func AddConfigFileFlag(flagSet *flag.FlagSet) {
	flagSet.String(configFile, "", "Configuration file in JSON, TOML, YAML, HCL, or Java properties formats (default none). See spf13/viper for precedence.")
}

// TryLoadConfigFile initializes viper with config file specified as flag
func TryLoadConfigFile(v *viper.Viper) error {
	if file := v.GetString(configFile); file != "" {
		v.SetConfigFile(file)
		err := v.ReadInConfig()
		if err != nil {
			return fmt.Errorf("cannot load config file %s: %w", file, err)
		}
	}
	return nil
}

// SharedFlags holds flags configuration
type SharedFlags struct {
	// Logging holds logging configuration
	Logging logging
}

type logging struct {
	Level    string
	Encoding string
}

// AddFlags adds flags for SharedFlags
func AddFlags(flagSet *flag.FlagSet) {
	AddLoggingFlags(flagSet)
}

// AddLoggingFlag adds logging flag for SharedFlags
func AddLoggingFlags(flagSet *flag.FlagSet) {
	flagSet.String(logLevel, "info", "Minimal allowed log Level. For more levels see https://github.com/uber-go/zap")
	flagSet.String(logEncoding, "json", "Log encoding. Supported values are 'json' and 'console'.")
}

// InitFromViper initializes SharedFlags with properties from viper
func (flags *SharedFlags) InitFromViper(v *viper.Viper) *SharedFlags {
	flags.Logging.Level = v.GetString(logLevel)
	flags.Logging.Encoding = v.GetString(logEncoding)
	return flags
}

// NewLogger returns logger based on configuration in SharedFlags
func (flags *SharedFlags) NewLogger(conf zap.Config, options ...zap.Option) (*zap.Logger, error) {
	var level zapcore.Level
	err := (&level).UnmarshalText([]byte(flags.Logging.Level))
	if err != nil {
		return nil, err
	}
	conf.Level = zap.NewAtomicLevelAt(level)
	conf.Encoding = flags.Logging.Encoding
	if flags.Logging.Encoding == "console" {
		conf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}
	return conf.Build(options...)
}
