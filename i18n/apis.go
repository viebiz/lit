package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	pkgerrors "github.com/pkg/errors"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func NewTranslator(supportedLang []string, basePath string) (Translator, error) {
	basePath = filepath.Clean(basePath)
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return Translator{}, fmt.Errorf("message path does not exist")
	}

	// Load language bundle
	b := i18n.NewBundle(defaultLangTag)
	b.RegisterUnmarshalFunc(defaultMessageFileFormat, json.Unmarshal)

	// Load all message file
	for _, lang := range supportedLang {
		msgPath := filepath.Join(basePath, fmt.Sprintf("%s.%s", lang, defaultMessageFileFormat))
		if _, err := b.LoadMessageFile(msgPath); err != nil {
			return Translator{}, pkgerrors.WithStack(err)
		}
	}

	t := Translator{
		defaultLang:     supportedLang[0], // The first item in supported language
		langLocalizeMap: make(map[string]*i18n.Localizer),
	}

	// Init localize map
	for _, lang := range supportedLang {
		t.langLocalizeMap[lang] = i18n.NewLocalizer(b, lang)
	}

	return t, nil
}
