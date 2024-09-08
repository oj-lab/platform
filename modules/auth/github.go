package auth_module

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	config_module "github.com/oj-lab/platform/modules/config"
)

const (
	githubOauthEntryURL  = "https://github.com/login/oauth/authorize"
	githubAccessTokenURL = "https://github.com/login/oauth/access_token"
	githubApiUserURL     = "https://api.github.com/user"

	servicePortProp       = "service.port"
	serviceHostProp       = "service.host"
	serviceSSLEnabledProp = "service.ssl_enabled"
)

var (
	githubClientID     string
	githubClientSecret string

	servicePort       uint
	serviceHost       string
	serviceSSLEnabled bool
)

func init() {
	githubClientID = config_module.AppConfig().GetString("auth_modulegithub.client_id")
	githubClientSecret = config_module.AppConfig().GetString("auth_modulegithub.client_secret")
	servicePort = config_module.AppConfig().GetUint(servicePortProp)
	serviceHost = config_module.AppConfig().GetString(serviceHostProp)
	serviceSSLEnabled = config_module.AppConfig().GetBool(serviceSSLEnabledProp)
}

func isGithubAuthEnabled() bool {
	return githubClientID != "" && githubClientSecret != ""
}

func GetGithubOauthEntryURL(callbackURL string) (*url.URL, error) {
	if !isGithubAuthEnabled() {
		return nil, fmt.Errorf("github auth is not enabled")
	}

	redirectUrl := fmt.Sprintf("%s:%d%s", serviceHost, servicePort, callbackURL)
	if serviceSSLEnabled {
		redirectUrl = "https://" + redirectUrl
	} else {
		redirectUrl = "http://" + redirectUrl
	}

	u, err := url.Parse(githubOauthEntryURL)
	if err != nil {
		return nil, err
	}

	query := u.Query()
	query.Add("client_id", githubClientID)
	query.Add("redirect_uri", redirectUrl)
	u.RawQuery = query.Encode()

	return u, nil
}

type GithubAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func GetGithubAccessToken(code string) (*GithubAccessTokenResponse, error) {
	u, err := url.Parse(githubAccessTokenURL)
	if err != nil {
		return nil, err
	}

	query := u.Query()
	query.Add("client_id", githubClientID)
	query.Add("client_secret", githubClientSecret)
	query.Add("code", code)

	u.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tokenResponse GithubAccessTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}

type GithubUser struct {
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

func GetGithubUser(accessToken string) (*GithubUser, error) {
	req, err := http.NewRequest(http.MethodGet, githubApiUserURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user GithubUser
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
