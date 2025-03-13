package i18n

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	pkgerrors "github.com/pkg/errors"
)

type localizer struct {
	underlyingLocalizer *i18n.Localizer
	getLocalizedMessage func(localizer *i18n.Localizer, cfg *i18n.LocalizeConfig) (string, error)
}

func newLocalizer(langKey string, bundle *i18n.Bundle) Localizable {
	return &localizer{
		underlyingLocalizer: i18n.NewLocalizer(bundle, langKey),
		getLocalizedMessage: getLocalizedMessage,
	}
}

func (l localizer) Localize(messageID string, params map[string]interface{}) string {
	msg, err := l.TryLocalize(messageID, params)
	if err != nil {
		return messageID
	}

	return msg
}

func (l localizer) TryLocalize(messageID string, params map[string]interface{}) (string, error) {
	if l.underlyingLocalizer == nil {
		return messageID, nil
	}

	// Prepare localize configuration
	cfg := &i18n.LocalizeConfig{
		MessageID: messageID,
	}
	if len(params) > 0 {
		cfg.TemplateData = params
	}

	// Localizable message
	rs, err := l.getLocalizedMessage(l.underlyingLocalizer, cfg)
	if err != nil {
		return "", pkgerrors.WithStack(err)
	}

	return rs, nil
}

func getLocalizedMessage(localizer *i18n.Localizer, cfg *i18n.LocalizeConfig) (string, error) {
	return localizer.Localize(cfg)
}
