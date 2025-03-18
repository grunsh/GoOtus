package settings

import "github.com/spf13/viper"

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

type LogSettings struct {
	Loglevel       LogLevel
	LoggingConsole bool
	LoggingFile    bool
	LogFileName    string
}

type Settings struct {
	LogSettings LogSettings
}

func GetSettings() Settings {
	f := viper.New()
	_ = f
	return Settings{}
}
