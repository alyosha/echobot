# echobot

As a basic example of how interactive message lifecycles work on Slack, this skeleton app responds with a message containing an interactive menu and button. Out of the box it simply allows users to select "participants" for a fictional event, and there is no validation to ensure a member is not selected more than once, etc.

Read more about interactive messages here: https://api.slack.com/interactive-messages

Feedback/suggestions welcome. 

# Setup

Getting the bot up and running is simple and the steps are as described below: 

1. After you've cloned the app, visit https://api.slack.com/apps and select `Create New App`
2. Once the app is created, jump inside and select `Add features and functionality`
3. Select `Interactive Components` and turn on Interactivity for your app
4. You will be prompted for a request URL: this is the API endpoint to which Slack will forward any messages your bot receives. While in the development process I highly recommend using ngrok, which can be downloaded at the following link (https://ngrok.com/https://ngrok.com/). Once you have your URL, append `/callback` and input on the `Interactive Components` page
5. Next select `Bot Users` from the sidebar and add a bot user to your app
6. Visit the `Oauth & Permissions` tab and install your bot user to your workspace. This will generate access tokens for both your app and the bot user
7. Obtain your signing secret from the `Basic Information` tab and set it as an environment variable (required)
8. Obtain your bot access token from the `Oauth & Permissions` page and set it as an environment variable (required)
9. Build the binary and start your server
10. Try to mention the bot from Slack and invite them to the channel of your choice when prompted
11. The `BOT_ID` environment variable is not set yet, so any subsequent attempts to mention the bot will simply result in the message contents being dumped in the logs. Mention the bot and check the log output to obtain the bot's ID
12. Set the bot's ID as an environment variable

You're done! The bot is now installed to your workplace and capable of responding to any mentions from channels you invite it to. Build on/change the existing response logic to create your own personal interactive message flow. 

If you need to persist data between requests, consider using a simple key-value cache such as https://github.com/patrickmn/go-cache and using the message timestamp as a key.
