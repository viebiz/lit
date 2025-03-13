package i18n

// MessageBundle represents the interface for localizing messages that support multiple languages
type MessageBundle interface {
	LoadMessageFile(path, langKey, ext string) error

	GetLocalize(langKey string) Localizable
}

// Localizable represents the interface for localizing messages
type Localizable interface {
	TryLocalize(messageID string, params map[string]interface{}) (string, error)

	Localize(messageID string, params map[string]interface{}) string
}
