# Web Based Twitch.tv Analytics Tool (Work In Progress)

TwitchSpy is a Go & Postgres application which logs the channels' statistics to monitor
their respective real-time popularity. Useful for streamers that would like to know
which games are increasing in popularity and which ones are declining.

The package comes with a minimal Go server to be used for GUI. The server supports encryption over HTTPS using
the Go autocert library by letsencrypt.org. The server is still under development.

## Currently Supported Features

### Server
TLS/SSL Support over port :443
HTTP over port :80
Debug/Reverse Proxy (Nginx, Apache etc) using any custom port

### Twitch
Real-Time Stream Data
Get Top Games
Get Top Channels for a Specific Game
Rich Game Details using GiantBomb API

## Planned Features

### Server
Wildcard TLS/SSL Support (delayed due to delayed release by letsencrypt.org)

### Twitch
Get Viewers
Monitor Viewer's Activity Per Stream
Multi-channel data gathering

## Config
The config is populated through the environment variables.