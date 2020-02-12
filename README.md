# AppEngine+Go Slack BOT template

## How to setup

1. https://api.slack.com/apps > `Create New App`
1. `Basic Information` > `Add features and functionality` > `Bots` > `Review Scopes to Add`
  * `app_mentions:read`, `chat:write`
  * `Install App to Workspace`
1. setup environment variables
  * `SLACK_BOT_VERIFICATION_TOKEN`
    * `Basic Information` > `Verification Token`
  * `SLACK_BOT_OAUTH_ACCESS_TOKEN`
    * `Install App` > `Bot User OAuth Access Token`
1. Deploy app to AppEngine
  * check `./deploy.sh`
1. `Basic Information` > `Add features and functionality` > `Event Subscriptions`
  * `Request URL`
  * `Subscribe to bot events` > `Add Bot User Event`
    * `app_mention`
  * `Save Changes`
