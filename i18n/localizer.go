package i18n

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	pkgerrors "github.com/pkg/errors"
)

var (
	prepareMessageConfigWrapper = prepareMessageConfig
)

type messageLocalize struct {
	localizer *i18n.Localizer
}

func (ml messageLocalize) Localize(messageID string, params map[string]interface{}) (string, error) {
	if ml.localizer == nil {
		return "", ErrBundleNotInitialized
	}

	rs, err := ml.localizer.Localize(prepareMessageConfigWrapper(messageID, params))
	if err != nil {
		return "", pkgerrors.WithStack(err)
	}

	return rs, nil
}

func prepareMessageConfig(messageID string, params map[string]interface{}) *i18n.LocalizeConfig {
	msgCfg := &i18n.LocalizeConfig{
		MessageID: messageID,
	}
	if len(params) > 0 {
		msgCfg.TemplateData = params
	}

	return msgCfg
}
