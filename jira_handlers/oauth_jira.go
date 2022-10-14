package jira_handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
)

type AccessibleResources []struct {
	Id string
}

var jiraOauthConfig = &oauth2.Config{
	RedirectURL:  os.Getenv("JIRA_URL_CALLBACK"),
	ClientID:     os.Getenv("JIRA_AUTH_APP_CLIENT_ID"),
	ClientSecret: os.Getenv("JIRA_AUTH_APP_CLIENT_SECRET"),
	Scopes:       []string{"read:me", "read:jira-work", "write:jira-work", "read:jira-user"},
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://auth.atlassian.com/authorize",
		TokenURL: "https://auth.atlassian.com/oauth/token",
	},
}

const jiraOauthURL = "https://api.atlassian.com/ex/jira/%s/%s"

func oauthJiraLogin(w http.ResponseWriter, r *http.Request) {

	// Create oauthState cookie
	oauthState := generateStateOauthCookie(w)

	/*
		AuthCodeURL receive state that is a token to protect the user from CSRF attacks. You must always provide a non-empty string and
		validate that it matches the the state query parameter on your redirect callback.
	*/
	u := jiraOauthConfig.AuthCodeURL(
		oauthState,
		oauth2.SetAuthURLParam("audience", "api.atlassian.com"),
		oauth2.SetAuthURLParam("state", "1"),
		oauth2.SetAuthURLParam("response_type", "code"),
		oauth2.SetAuthURLParam("prompt", "consent"),
	)

	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func oauthJiraCallback(w http.ResponseWriter, r *http.Request) {
	// Read oauthState from Cookie
	// oauthState, _ := r.Cookie("oauthstate")

	// if r.FormValue("state") != oauthState.Value {
	// 	log.Println("invalid oauth google state")
	// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	// 	return
	// }

	data, err := getUserDataFromJira(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// GetOrCreate User in your db.
	// Redirect or response with a token.
	// More code .....
	fmt.Fprintf(w, "UserInfo: %s\n", data)
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(20 * time.Minute)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func getUserDataFromJira(code string) ([]byte, error) {
	// Use code to get token and get user info from Google.
	token, err := jiraOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}

	// GET APP ID
	request, err := http.NewRequest(
		"GET",
		"https://api.atlassian.com/oauth/token/accessible-resources",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed getting request")
	}

	ctx := context.Background()
	tokenSource := jiraOauthConfig.TokenSource(ctx, token)
	client := oauth2.NewClient(ctx, tokenSource)

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed getting accessible resources: %s", err.Error())
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read accesible resources response: %s", err.Error())
	}

	var resp AccessibleResources
	err = json.Unmarshal(contents, &resp)

	if err != nil {
		return nil, fmt.Errorf("failed marshall: %s", err.Error())
	}

	request, err = http.NewRequest(
		"GET",
		fmt.Sprintf(jiraOauthURL, resp[0].Id, "rest/api/2/myself"),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed getting request myself req")
	}

	response, err = client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer response.Body.Close()
	contents, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}

	return contents, nil
}
