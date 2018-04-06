# Web Based Twitch.tv Analytics Tool (Work In Progress)

TwitchSpy is a Go & Postgres application which logs the channels' statistics to monitor
their respective real-time popularity. Useful for streamers that would like to know
which games are increasing in popularity and which ones are declining.

The package comes with a minimal Go server to be used for GUI. The server supports encryption over HTTPS using
the Go autocert library by letsencrypt.org. The server is still under development.

## Currently Supported Features

### Server
1. TLS/SSL Support over port :443
2. HTTP over port :80
3. Debug/Reverse Proxy (Nginx, Apache etc) using any custom port

### Twitch
1. Real-Time Stream Data
2. Get Top Games
3. Get Top Channels for a Specific Game
4. Rich Game Details using GiantBomb API
5. Database storage

## Planned Features

### Server
Wildcard TLS/SSL Support (delayed due to delayed release by letsencrypt.org)

### Twitch
1. Real-Time Viewers' Data
2. Monitor Viewer's Activity Per Stream
3. Multi-channel data gathering

## Config
The config is populated through the environment variables.