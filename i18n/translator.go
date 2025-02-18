package i18n

import (
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	pkgerrors "github.com/pkg/errors"
	"golang.org/x/text/language"
)

var (
	defaultLangTag = language.English
)

const (
	defaultMessageFileFormat = "json"
)

type Translator struct {
	defaultLang     string
	langLocalizeMap map[string]*i18n.Localizer
}

func (tr Translator) Translate(msgID string, params map[string]any) (string, error) {
	return tr.TranslateWithLang(tr.defaultLang, msgID, params)
}

func (tr Translator) TranslateWithLang(lang string, msgID string, params map[string]any) (string, error) {
	localizer, exist := tr.langLocalizeMap[lang]
	if !exist {
		return "", fmt.Errorf("%w: %s", ErrGivenLangNotSupported, lang)
	}

	localizeCfg := i18n.LocalizeConfig{
		MessageID: msgID,
	}
	if len(params) > 0 {
		localizeCfg.TemplateData = params
	}

	result, err := localizer.Localize(&localizeCfg)
	if err != nil {
		return "", pkgerrors.WithStack(err)
	}

	return result, nil
}
