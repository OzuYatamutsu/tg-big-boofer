# tg-big-boofer [![Build Status](https://travis-ci.org/OzuYatamutsu/tg-big-boofer.svg?branch=master)](https://travis-ci.org/OzuYatamutsu/tg-big-boofer)
A Telegram bot to guard against low-effort spambots in large groups.
<div align="center">
    <img src="https://raw.githubusercontent.com/OzuYatamutsu/tg-big-boofer/master/bigboofer.png" width="300" height="300" /><br />
</div>

## To run
BigBoofer uses go modules, so it requires Go >= 1.11.0.

Contact `@BotFather` on Telegram for an API key. Then, set your API key in `API_KEY.config` and run:

```
go build
```

It should pull in all required dependencies and produce a binary ready for you to run.

## To test
```
go test -v bigboofer/test
```