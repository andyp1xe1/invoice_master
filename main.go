package main

import "log/slog"

func main() {
	//llama("json example of an invoice. only json. fast")
	server, err := NewServer("localhost:1337")
	if err != nil {
		return
	}
	if err := server.Run(); err != nil {
		slog.Error("FATAL:" + " " + err.Error())
		return
	}
	slog.Info("Exiting ...")
}
