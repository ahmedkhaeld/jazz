package urlsigner

import (
	"fmt"
	goalone "github.com/bwmarrin/go-alone"
	"strings"
	"time"
)

type Signer struct {
	Secret []byte
}

// GenerateSignedURL appends a hash to the provided url and that creates a signedUrl that's unique
func (s *Signer) GenerateSignedURL(url string) string {
	var urlToSign string
	crypt := goalone.New(s.Secret, goalone.Timestamp)

	if strings.Contains(url, "?") {
		urlToSign = fmt.Sprintf("%s&hash=", url)
	}
	if !strings.Contains(url, "?") {
		urlToSign = fmt.Sprintf("%s?hash=", url)
	}
	urlBytes := crypt.Sign([]byte(urlToSign))
	signedURL := string(urlBytes)

	return signedURL

}

func (s *Signer) Verify(signedUrl string) bool {
	crypt := goalone.New(s.Secret, goalone.Timestamp)
	_, err := crypt.Unsign([]byte(signedUrl))
	if err != nil {
		return false
	}
	return true
}

// IsExpired checks if a signedUrl has or not expired yet
func (s *Signer) IsExpired(signedUrl string, minutesUntilExpire int) bool {
	crypt := goalone.New(s.Secret, goalone.Timestamp)
	ts := crypt.Parse([]byte(signedUrl))
	return time.Since(ts.Timestamp) > time.Duration(minutesUntilExpire)*time.Minute
}
