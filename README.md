# tg-big-boofer [![Build Status](https://travis-ci.org/OzuYatamutsu/tg-big-boofer.svg?branch=master)](https://travis-ci.org/OzuYatamutsu/tg-big-boofer)
A Telegram bot to guard against low-effort spambots in large groups.
<div align="center">
    <img src="https://raw.githubusercontent.com/OzuYatamutsu/tg-big-boofer/master/bigboofer.png" width="300" height="300" /><br />
</div>

## To add to your group
Add `@BigBooferBot` to your group. For enforcement to work (see below),
they must be an admin of the group they are a part of.

Once added to the group, if you are an admin, mention `@BigBooferBot` 
in your group with the link to the channel and passphrase, and it will 
begin enforcement.

## When a new user joins your group...
* They will be welcomed by **a big friend,** `@BigBooferBot`.
* They will be directed to a separate channel, containing a passphrase 
(and whatever else you want, e.g. rules).
* Until they reply in the channel with the passphrase, all of their messages 
will be deleted as soon as they are posted.
* If they don't reply with the passphrase within 5 minutes, `@BigBooferBot` 
will (regretably) remove them from the group.

## To run
BigBoofer uses go modules, so it requires Go >= 1.11.0.

Contact `@BotFather` on Telegram for an API key. Then, set your API key 
in `config.go` and run:

```
go build
```

It should pull in all required dependencies and produce a binary ready 
for you to run.

## To test
```
go test -v bigboofer/test
```