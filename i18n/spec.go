package i18n

// Bundle represents the interface for localizing messages that support multiple languages
type Bundle interface {
	LoadMessageFile(path, langKey, ext string) error
	
	Localize(messageID string, params map[string]interface{}) (string, error)

	LocalizeWithLang(langKey string, messageID string, params map[string]interface{}) (string, error)

	TryLocalize(messageID string, params map[string]interface{}) string
}

// MessageLocalize represents the interface for localizing messages
type MessageLocalize interface {
	Localize(messageID string, params map[string]interface{}) (string, error)
}
