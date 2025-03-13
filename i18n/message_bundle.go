package i18n

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	pkgerrors "github.com/pkg/errors"
	"github.com/viebiz/lit/monitoring"
)

type bundleMessage struct {
	DefaultLang      string
	SourcePath       string
	BundleFileFormat string
	underlyingBundle *i18n.Bundle
	monitor          *monitoring.Monitor
	langOnce         sync.Map
	localizeMap      map[string]Localizable // For caching localizer
}

func (b *bundleMessage) LoadMessageFile(path string, langKey, ext string) error {
	if b == nil {
		return nil
	}

	return b.loadMessageFile(path, langKey, ext)
}

func (b *bundleMessage) loadMessageFile(path string, langKey string, ext string) error {
	// Validate source path
	if err := validateSourcePath(path); err != nil {
		return pkgerrors.WithStack(err)
	}

	// Load message file
	msgPath := filepath.Join(path, fmt.Sprintf("%s.%s", langKey, ext))
	if _, err := b.underlyingBundle.LoadMessageFile(msgPath); err != nil {
		return pkgerrors.WithStack(err)
	}

	// Register localize
	b.localizeMap[langKey] = newLocalizer(langKey, b.underlyingBundle)

	return nil
}

func (b *bundleMessage) GetLocalize(langKey string) Localizable {
	if lc, exists := b.localizeMap[langKey]; exists {
		return lc
	}

	// Ensure a sync.Once is set for the given langKey
	onceInterface, _ := b.langOnce.LoadOrStore(langKey, &sync.Once{})
	once := onceInterface.(*sync.Once)

	// Load message file once for each language key
	once.Do(func() {
		if err := b.LoadMessageFile(b.SourcePath, langKey, b.BundleFileFormat); err != nil {
			b.monitor.Errorf(err, "Failed to load message file")
		}
	})

	return b.localizeMap[langKey]
}
