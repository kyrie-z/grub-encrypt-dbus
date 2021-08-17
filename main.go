package main

import "grub-encrypt-dbus/pkg/server"

func main() {
	srv := server.GetService()
	err := srv.Init()
	if err != nil {
		panic(err)
	}

	srv.Loop()
}
