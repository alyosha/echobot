# echobot

A basic example of interactive message flows in Slack, this skeleton app responds to a slash command invocation with an action block containing a select menu and button. Out of the box, it simply allows users to select/deselect members from their workspace.

The purpose of this repo is to reduce boilerplate when starting up a new Slackbot project. Build on/change the existing response logic to create your own personal interactive message flow. 

Blocks are Slack's replacement for the [now-deprecated](https://api.slack.com/messaging/attachments-to-blocks) message attachments. Read more about them [here](https://api.slack.com/reference/messaging/blocks).


## Setup
Follow the steps below to get the bot up and running: 

1. After you've cloned the app, visit the Slack API dashboard and select `Create New App`.
2. Once the app is created, jump inside and select `Add features and functionality`.
3. Select `Interactive Components` and turn on interactivity for your app.
4. You will be prompted for a request URL: this is the API endpoint to which Slack will forward any callback messages your bot receives. While in the development process I highly recommend using [ngrok](https://ngrok.com/). Once you have your URL, append `/callback` and input it on the `Interactive Components` page.
5. Setup two slash commands for the `/add-users` and `/help` endpoints from the
   Slash Commands tab of the dashboard
6. Next, select `Bot Users` from the sidebar and add a bot user to your app.
7. Visit the `Oauth & Permissions` tab and install the app to your workspace. This will generate access tokens for both you and the bot user.
8. Obtain your signing secret from the `Basic Information` tab and set it as an environment variable (required).
9. Obtain your bot's access token from the `Oauth & Permissions` page and set it as an environment variable (required).
10. Build the binary and start your server.
    - This project uses modules, which were first introduced in Go `1.11`. Use
   of modules is optional in `1.11`, so if you are working within your `GOPATH`
you will need to enable the use of modules by running: `export GO111MODULE=on`

You're done! The bot is now installed to your workplace and capable of responding to the registered slash commands from any channels you invite it to.

## Deployment
The project includes a basic deploy script for use with Google's App Engine.
You will need to have the [gcloud CLI](https://cloud.google.com/sdk/gcloud/) 
installed to utilize this locally, otherwise you can deploy via the [cloud shell](https://cloud.google.com/shell/docs/).

By default, the app will be deployed to the GCP project specified via `gcloud config set core/project PROJECT`. 

You can confirm your current project by running the following command: `gcloud config get-value project` 
