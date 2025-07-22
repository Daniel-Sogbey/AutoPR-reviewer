package githubapi

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"os"
	"review-pr/webhook-service/internal/github"
	"review-pr/webhook-service/internal/requester"
	"time"
)

func loadPrivateKey(keyPath string) (*rsa.PrivateKey, error) {
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("invalid PEM format")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func generateJWT(appID string, privateKey *rsa.PrivateKey) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"iat": now.Add(-time.Minute).Unix(),
		"exp": now.Add(time.Minute * 9).Unix(),
		"iss": appID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(privateKey)
}

func Auth() (*github.InstallationTokenResponseModel, error) {
	privateKey, err := loadPrivateKey(os.Getenv("GITHUB_APP_PRIVATE_KEY_PATH"))
	
	if err != nil {
		return nil, err
	}

	jwtToken, err := generateJWT("1472248", privateKey)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Add("Accept", "application/vnd.github+json")
	headers.Add("X-GitHub-Api-Version", "2022-11-28")
	installationResponse, err := requester.Requester[github.InstallationsResponseModel](http.MethodGet, "https://api.github.com/repos/Daniel-Sogbey/code-reviewer/installation", jwtToken, headers, nil)
	if err != nil {
		return nil, err
	}

	headers = http.Header{}
	headers.Add("Accept", "application/vnd.github+json")
	headers.Add("X-GitHub-Api-Version", "2022-11-28")
	installationTokenResponse, err := requester.Requester[github.InstallationTokenResponseModel](http.MethodPost, fmt.Sprintf("https://api.github.com/app/installations/%d/access_tokens", installationResponse.ID), jwtToken, headers, nil)
	if err != nil {
		return nil, err
	}

	log.Println(installationResponse)
	return installationTokenResponse, nil
}
