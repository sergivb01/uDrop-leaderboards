package main

import (
	service "github.com/sergivb01/udrop-leaderboards/services"
)

func main() {
	svc := service.NewService("blablabla")

	done := make(chan bool, 1)

	svc.ListenHTTP(done)
	<-done
	svc.Logger.Println("server has stopped")
}
