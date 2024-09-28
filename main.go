package main

import "log/slog"

func main() {
	llama("json example of an invoice. only json. fast")
	return
	server := NewServer("localhost:1337")
	if err := server.Run(); err != nil {
		slog.Error("FATAL:" + " " + err.Error())
		return
	}
	slog.Info("Exiting ...")
}
