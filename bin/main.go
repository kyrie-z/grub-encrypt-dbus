package main

import (
	"grub-encrypt-dbus/grub2/server"

	"pkg.deepin.io/lib/log"
)

var logger = log.NewLogger("zzl/grub")

func init() {
	server.SetLogger(logger)

}

func main() {
	srv := server.GetService()
	err := srv.Init()
	if err != nil {
		panic(err)
	}
	logger.Debug("mode:debug")
	logger.Info("mode:info")
	srv.Loop()
}
