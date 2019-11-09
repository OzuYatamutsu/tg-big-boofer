# tg-big-boofer [![Build Status](https://travis-ci.org/OzuYatamutsu/tg-big-boofer.svg?branch=master)](https://travis-ci.org/OzuYatamutsu/tg-big-boofer)
A Telegram bot to guard against low-effort spambots in large groups.
<div align="center">
    <img src="https://raw.githubusercontent.com/OzuYatamutsu/tg-big-boofer/master/bigboofer.png" width="300" height="300" /><br />
</div>

## To add to your group
Add `@BigBooferBot` to your group. For enforcement to work (see below),
they must be an admin of the group they are a part of.

![a friend](https://raw.githubusercontent.com/OzuYatamutsu/tg-big-boofer/master/bigboofer-demo00.png)

Once added to the group, if you are an admin, promote `@BigBooferBot` 
to an admin, and configure the passphrase via `/setchannel <channel_url> <passphrase>`:

![a friend](https://raw.githubusercontent.com/OzuYatamutsu/tg-big-boofer/master/bigboofer-demo05.png)

`@BigBooferBot` will then begin enforcement.

## When a new user joins your group...
* They will be welcomed by **a big friend,** `@BigBooferBot`, and will be
directed to a separate channel, containing a passphrase (and whatever else
you want, e.g. rules).
![a friend](https://raw.githubusercontent.com/OzuYatamutsu/tg-big-boofer/master/bigboofer-demo01.png)

* Until they reply in the channel with the passphrase, all of their messages will either not
be allowed to be posted, or deleted as soon as they are posted.
![a friend](https://raw.githubusercontent.com/OzuYatamutsu/tg-big-boofer/master/bigboofer-demo02.png)

* If they don't reply with the passphrase within 5 minutes, `@BigBooferBot` 
will (regretably) remove them from the group.
![a friend](https://raw.githubusercontent.com/OzuYatamutsu/tg-big-boofer/master/bigboofer-demo03.png)

* ...but admins can manually approve new users at any time.
![a friend](https://raw.githubusercontent.com/OzuYatamutsu/tg-big-boofer/master/bigboofer-demo04.png)

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