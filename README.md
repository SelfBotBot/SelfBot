# SelfBot
A Discord soundboard bot, written in Go.

## Information
The ultimate soundboard bot with plenty of planned features.

SelfBot is simple to use, easy to add custom sounds to, ~~with a simple and ergonomic web interface.~~ soon:tm:

## TODO
- [ ] Sound management and storage not in fs.
- [ ] Finish web UI (start v2)
- [ ] Add proper support for sharding
- [ ] Redo audio handler (maybe); use locks and what not instead of channels.
- [ ] Write apps and programs for hotkey/mobile phone support.
- [ ] Alexa integration
- [ ] Google assistant integration

## Self hosting
~~Not sure how well this will work, some code might need changing due to SSL and stuff..~~


### Requirments
- Golang
- FFmpeg
- Redis
- MariaDB/MySQL

### Execution/build
The easiest way to run this program is to use the golang path and git, that way you won't have to worry about all of the assets and audio file locations.
1. `git clone github.com/selfbotbot/selfbot`
2. `go mod download`
3. `go run cmd/selfbot-bot/main.go`

A configuration file will be created and the program will exit and ask for it to be edited. Fill in the apropriate information and restart the bot.

