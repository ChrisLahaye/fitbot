# fitbot

This package provides a Telegram bot to automatically make reservations at Fit For Free.

## Build

```sh
git clone https://github.com/ChrisLahaye/fitbot.git
go build -o fitbot cmd/main.go
```

## Usage

```sh
./fitbot
```

## Configuration
The following variables are all required and must be specified as environment variables or in the `.env` configuration file.

### `BOT_OWNER`

The Telegram user ID communicating with the bot.

You can view your Telegram user ID using the [GetIDs](https://t.me/getidsbot) bot.

### `BOT_TOKEN`

The Telegram API token of the bot.

You can view and manage your Telegram bots using the [BotFather](https://t.me/botfather) bot.

### `FIT_MEMBER_ID`

The Fit For Free membership number.

### `FIT_POSTCODE`

The postcode of the Fit Fot Free membership.

### `FIT_VENUE`

The Fit For Free venue to make reservations at.

You can list all venues by querying the `/v1/venues/` endpoint of the Fit For Free API. Check the source code for how to make authorized API requests.

---

*The Fit For Free API requires an app version to be specified as request header. Outdated app versions are blocked and result in an error message requesting to update the app. This request header is hard-coded [here](https://github.com/ChrisLahaye/fitbot/blob/develop/internal/fitforfree/api.go#L26) with the latest Android app version at the time of development. To successfully use the API this value must be a recent app version. Check the [Google Play](https://play.google.com/store/apps/details?id=nl.fitforfree.serviceapp) store for the latest app version.*

*For Educational Purposes Only*
