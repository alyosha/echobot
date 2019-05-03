# echobot

A basic example of how Slack's interactive messaging works, this skeleton app responds to self-mentions with a message containing an interactive menu and button. Out of the box it simply allows users to select "participants" for a fictional event, and there is no validation to ensure a member is not selected more than once, etc.

The purpose of this repo is to reduce boilerplate when starting up a new Slackbot project. Build on/change the existing response logic to create your own personal interactive message flow. 

In my projects I always configure help flows via slash command. I find that slash commands are more helpful than RTM in this regard in that, if the user forgets the help command itself, they can easily remind themselves by looking at the bot's profile page. `echobot` comes with a pre-configured `help` endpoint which can be used after setting up a new slash command on the [Slack API dashboard](https://api.slack.com/apps).

Read more about interactive messages [here](https://api.slack.com/interactive-messages).


# Setup

Follow the steps below to get the bot up and running: 

1. After you've cloned the app, visit the Slack API dashboard and select `Create New App`.
2. Once the app is created, jump inside and select `Add features and functionality`.
3. Select `Interactive Components` and turn on interactivity for your app.
4. You will be prompted for a request URL: this is the API endpoint to which Slack will forward any callback messages your bot receives. While in the development process I highly recommend using [ngrok](https://ngrok.com/). Once you have your URL, append `/callback` and input it on the `Interactive Components` page.
5. Next, select `Bot Users` from the sidebar and add a bot user to your app.
6. Visit the `Oauth & Permissions` tab and install the app to your workspace. This will generate access tokens for both you and the bot user.
7. Obtain your signing secret from the `Basic Information` tab and set it as an environment variable (required).
8. Obtain your bot's access token from the `Oauth & Permissions` page and set it as an environment variable (required).
9. Build the binary and start your server.
10. Try to mention the bot in a channel of your choice and invite them to join when prompted.
11. In order to respond to RTM message events, `echobot` requires the `BOT_ID` environment variable to be set. Any messages received before the `BOT_ID` is configured will be dumped in the logs, so mention the bot and check your log output to obtain the bot's ID.
12. Set the bot's ID as an environment variable.

You're done! The bot is now installed to your workplace and capable of responding to any mentions from channels you invite it to.

If you need to persist data between requests, consider using a simple key-value cache such as https://github.com/patrickmn/go-cache and using the message timestamp as a key.
