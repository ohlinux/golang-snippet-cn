package logs

import (
	seelog "github.com/cihub/seelog"
)

var Logger seelog.LoggerInterface

func LoadConf(configPath string) {
	DisableLog()
	logger, err := seelog.LoggerFromConfigAsFile(configPath)
	if err != nil {
		panic(err)
	}
	UseLogger(logger)
}

// DisableLog disables all library log output
func DisableLog() {
	Logger = seelog.Disabled
}

// UseLogger uses a specified seelog.LoggerInterface to output library log.
// Use this func if you are using Seelog logging system in your app.
func UseLogger(newLogger seelog.LoggerInterface) {
	Logger = newLogger
}
