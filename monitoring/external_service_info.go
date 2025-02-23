package monitoring

import (
	"net/url"
)

// ExternalServiceInfo holds the ext svc info
type ExternalServiceInfo struct {
	Hostname string
	Port     string
}

// NewExternalServiceInfo creates a new ExternalServiceInfo from the given url
func NewExternalServiceInfo(rawURL string) ExternalServiceInfo {
	info := ExternalServiceInfo{}

	u, err := url.Parse(rawURL)
	if err != nil || u.Host == "" {
		u, err = url.Parse("https://" + rawURL)
		if err != nil {
			return info
		}
	}
	info.Hostname = u.Hostname()
	info.Port = u.Port()

	return info
}
