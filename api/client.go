package api

import (
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/spf13/viper"
)

var (
	Client = &http.Client{Timeout: 60 * time.Second}
)

func Url(route string) string {
	u, _ := url.Parse("https://api.fold.money/api")
	u.Path = path.Join(u.Path, route)
	return u.String()
}

func APIRequest(method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, path, body)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+viper.GetString("token.access"))
	req.Header.Set("User-Agent", "ua/unfold")
	req.Header.Set("X-Device-Hash", viper.GetString("device_hash"))
	req.Header.Set("X-Device-Location", "India")
	req.Header.Set("X-Device-Name", "unfold")
	req.Header.Set("X-Device-Type", "Android")

	return req, nil
}
