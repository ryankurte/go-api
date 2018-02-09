package plugins

import (
	log "github.com/sirupsen/logrus"
)

// LogPlugin plugin context structure
type LogPlugin struct {
	logger log.FieldLogger
}

// NewLogPlugin Create a new registration logging plugin instance
func NewLogPlugin() *LogPlugin {
	return &LogPlugin{
		logger: log.New().WithField("module", "log-plugin"),
	}
}

// Register handler called when new paths are registered
func (l *LogPlugin) Register(route string, method string, input interface{}, output interface{}) {
	l.logger.Infof("Registered handler for route: %s method: %s with input: %+v, output: %+t\n",
		route, method, input, output)
}
