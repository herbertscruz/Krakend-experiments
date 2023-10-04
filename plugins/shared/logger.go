package shared

import "fmt"

var errorTemplate string

type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warning(v ...interface{})
	Error(v ...interface{})
	Critical(v ...interface{})
	Fatal(v ...interface{})
}

type CustomLogger struct{}

func (n CustomLogger) Debug(_ ...interface{})    {}
func (n CustomLogger) Info(_ ...interface{})     {}
func (n CustomLogger) Warning(_ ...interface{})  {}
func (n CustomLogger) Error(_ ...interface{})    {}
func (n CustomLogger) Critical(_ ...interface{}) {}
func (n CustomLogger) Fatal(_ ...interface{})    {}

func template(level string, pluginName string, message string) string {
	return fmt.Sprintf("[%s][PLUGIN: %s] %s", level, pluginName, message)
}

func RegisterCustomLogger(logger Logger, pluginName string, v interface{}) {
	l, ok := v.(Logger)
	if !ok {
		return
	}
	logger = l
	logger.Debug(template("DEBUG", pluginName, "Logger loaded"))
}
