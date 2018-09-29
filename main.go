package main

import (
	"os"
	"os/signal"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func main() {
	parseConfig()

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		stopServer()
		initConfig()
		go startServer()
	})

	die := make(chan bool)
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	go func() {
		<-sig
		die <- true
	}()

	go startServer()
	<-die

	stopServer()
}
