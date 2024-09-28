package main

import "log/slog"

func main() {
	server := NewServer("localhost:1337")
	if err := server.Run(); err != nil {
		slog.Error("FATAL:" + " " + err.Error())
		return
	}
	slog.Info("Exiting ...")
}

