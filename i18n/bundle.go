package i18n

import (
	"fmt"
	"path/filepath"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	pkgerrors "github.com/pkg/errors"
)

type bundle struct {
	i18nBundle  *i18n.Bundle
	DefaultLang string
	LocalizeMap map[string]MessageLocalize
}

func (b *bundle) LoadMessageFile(path string, langKey, ext string) error {
	if b == nil {
		return nil
	}

	return b.loadMessageFile(path, langKey, ext)
}

func (b *bundle) loadMessageFile(path string, langKey string, ext string) error {
	// Validate source path
	if err := validateSourcePath(path); err != nil {
		return pkgerrors.WithStack(err)
	}

	// Load message file
	msgPath := filepath.Join(path, fmt.Sprintf("%s.%s", langKey, ext))
	if _, err := b.i18nBundle.LoadMessageFile(msgPath); err != nil {
		return pkgerrors.WithStack(err)
	}

	// Register localize
	b.LocalizeMap[langKey] = messageLocalize{
		localizer: i18n.NewLocalizer(b.i18nBundle, langKey),
	}

	return nil
}

func (b *bundle) LocalizeWithLang(langKey string, messageID string, params map[string]interface{}) (string, error) {
	// Skip if bundle is nil
	if b == nil {
		return messageID, nil
	}

	localize, exists := b.LocalizeMap[langKey]
	if !exists {
		return "", ErrGivenLangNotSupported
	}

	rs, err := localize.Localize(messageID, params)
	if err != nil {
		return "", err
	}

	return rs, nil
}

func (b *bundle) Localize(messageID string, params map[string]interface{}) (string, error) {
	return b.LocalizeWithLang(b.DefaultLang, messageID, params)
}

func (b *bundle) TryLocalize(messageID string, params map[string]interface{}) string {
	msg, err := b.Localize(messageID, params)
	if err != nil {
		return messageID
	}

	// Skip error and return messageID
	return msg
}
