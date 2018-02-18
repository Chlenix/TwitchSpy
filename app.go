package main

import "TwitchSpy/service"

func main() {
	Dispatcher := service.NewDispatcher()
	Dispatcher.Start()
}
