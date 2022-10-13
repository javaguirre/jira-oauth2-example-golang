# Oauth2-example with Go for JIRA user connection
Authentication is the most common part in any application. You can implement your own authentication system or use one of the many alternatives that exist, but in this case we are going to use OAuth2.

OAuth is a specification that allows users to delegate access to their data without sharing
their username and password with that service, if you want to read more about Oauth2 go [here](https://oauth.net/2/).


## Config JIRA
You need to create a JIRA app and set the OAuth callback URL in [the Atlassian dev console](https://developer.atlassian.com/console/myapps/XXXXXXXXX/authorization/auth-code-grant). You can get the client ID and Secret from the Jira App console settings.

Then export envvars:

```bash
export JIRA_URL_CALLBACK=<MYURL>/auth/jira/callback
export JIRA_AUTH_APP_CLIENT_ID=<MY_CLIENT_ID_FROM_JIRA_APP>
export JIRA_AUTH_APP_CLIENT_SECRET=<MY_CLIENT_SECRET_FROM_JIRA_APP>
```

# Running

```bash
## let's run and test
go run main.go
```
