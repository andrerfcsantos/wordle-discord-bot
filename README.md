# Wordle Discord Bot

Discord bot for the popular [Wordle](https://www.nytimes.com/games/wordle/index.html) game.

The bot scans discord channels for Wordle copy/pastes and keeps track of the score of each person in the channel.

## Using this bot in your server

* Invite the bot to your server using [this link](https://discord.com/api/oauth2/authorize?client_id=934221170897743882&permissions=534723951680&scope=applications.commands%20bot).
* Do `/wordle track` in the channel you want the bot to listen to wordle copy/pastes.

### Commands

* `/wordle track`: Start tracking wordle copy/pastes in the current channel.
* `/wordle leaderboard`: Prints the leaderboard of the current channel.
* More commands coming soon!

## Running the bot on your own server/machine

You can run the bot on your own server/machine by building it locally or using docker.

### Running locally without docker

This requires having Go installed and a postgres database:

* Clone the repository
  * `git clone git@github.com:andrerfcsantos/wordle-discord-bot.git`
* Build the bot
  * `cd wordle-discord-bot && go build`
* Create a bot application in Discord
  * You can do it in the [Discord Developers: Applications](https://discord.com/developers/applications/) page.
* Setup a postgres database to be used by the bot or use an existing one.
* Set the environment variables (see [below](#environment-variables))
* Run the bot
  * `./wordle-discord-bot`
* To test the bot, create a discord server and invite the bot to it.
  * To invite a bot to your server, paste this link in a browser, replacing `<client_id>` with the client id of your bot (also known as app id):
  `https://discord.com/api/oauth2/authorize?client_id=<client_id>&permissions=534723951680&scope=applications.commands%20bot`

### Running with docker

Running the bot with docker doesn't require a database previously configured or having Go installed.
The only requirement is having Docker installed:

* Clone the repository
  * `git clone git@github.com:andrerfcsantos/wordle-discord-bot.git`
* Create a bot application in Discord
  * You can do it in the [Discord Developers: Applications](https://discord.com/developers/applications/) page.
* Set the environment variables in `docker-compose.yml`
  * The variables you need to set are `WORDLE_DISCORD_BOT_APP_ID` and `WORDLE_DISCORD_BOT_TOKEN`
  * The other variables are already set for you, but feel free to change them as you see fit.
* Run the bot
  * `docker-compose up -d`
* To test the bot, create a discord server and invite the bot to it.
  * To invite a bot to your server, paste this link in a browser, replacing `<client_id>` with the client id of your bot (also known as app id):
    `https://discord.com/api/oauth2/authorize?client_id=<client_id>&permissions=534723951680&scope=applications.commands%20bot`

### Environment Variables

The bot requires the following environment variables to be set:

* `WORDLE_DISCORD_BOT_APP_ID`
  * The application ID of the bot.
  You can see this in the page for your application in the [Discord Developers Portal](https://discord.com/developers/applications).
* `WORDLE_DISCORD_BOT_TOKEN`
  * Bot token.
  This is used for authentication with discord.
  You can see this in the page for your application in the [Discord Developers Portal](https://discord.com/developers/applications).
* `WORDLE_DB_HOST`
  * Postgres database host (e.g `localhost`). 
* `WORDLE_DB_PORT`
  * Postgres database port (e.g `5432`). 
* `WORDLE_DB_USER`
  * Postgres database user (e.g `wordle`).
* `WORDLE_DB_PASSWORD`
  * Postgres database password for the `WORDLE_DB_USER` user (e.g `admin1234`).
* `WORDLE_DB_NAME`
  * Postgres database name (e.g `wordle`).

## Contributing

Contributions in the form of PRs are welcome.
