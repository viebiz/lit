package iam

import (
	"strings"

	"github.com/viebiz/lit/monitoring"
)

var (
	// EnforcerLog is exported variable to support log during development
	EnforcerLog = false
)

type enforcerLogger struct {
	*monitoring.Monitor
}

func (i enforcerLogger) EnableLog(b bool) {}

func (i enforcerLogger) IsEnabled() bool {
	return EnforcerLog
}

func (i enforcerLogger) LogModel(model [][]string) {
	i.Infof("Model %v", model)
}

func (i enforcerLogger) LogEnforce(matcher string, request []interface{}, result bool, explains [][]string) {
	if !i.IsEnabled() {
		return
	}

	i.Infof("Enforced %v, (%s), (%v), hit policy: %v", result, request, matcher, explains)
}

func (i enforcerLogger) LogRole(roles []string) {
	if !i.IsEnabled() {
		return
	}

	i.Infof("Roles: %s", roles)
}

func (i enforcerLogger) LogPolicy(policy map[string][][]string) {
	if !i.IsEnabled() {
		return
	}

	i.Infof("Policies: %v", policy)
}

func (i enforcerLogger) LogError(err error, msg ...string) {
	if !i.IsEnabled() {
		return
	}

	i.Errorf(err, strings.Join(msg, ","))
}
