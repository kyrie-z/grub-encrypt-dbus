package server

import (
	"fmt"
	"grub-encrypt-dbus/pkg/utils"

	"github.com/godbus/dbus"
	"pkg.deepin.io/lib/dbusutil"
)

func (g *GrubEncrypt) AddAccount(usr, passwd string) *dbus.Error {
	chiperPasswd, err := utils.GrubPBKDF2Crypto(usr, passwd)
	if err != nil {
		fmt.Println(err)
		return dbusutil.ToError(err)
	}
	fmt.Println(chiperPasswd)
	err = utils.Add.WriteConfig(usr, chiperPasswd)
	if err != nil {
		fmt.Println(err)
		return dbusutil.ToError(err)
	}
	g.Status = "enable"
	g.User = append(g.User, usr)
	return nil
}

func (g *GrubEncrypt) DeleteAccount(usr string) *dbus.Error {
	err := utils.Delete.WriteConfig(usr)
	if err != nil {
		fmt.Println(err)
		return dbusutil.ToError(err)
	}
	// g.Status = "disable"
	return nil
}

func (g *GrubEncrypt) DisableAuthentication(usr string) *dbus.Error {
	err := utils.Disable.WriteConfig(usr)
	if err != nil {
		fmt.Println(err)
		return dbusutil.ToError(err)
	}
	return nil
}

func (g *GrubEncrypt) EnableAuthentication(usr string) *dbus.Error {
	err := utils.Enable.WriteConfig(usr)
	if err != nil {
		fmt.Println(err)
		return dbusutil.ToError(err)
	}
	return nil
}

func (g *GrubEncrypt) detectStatus() {
	//检测是否有有效账户
	g.Status = "detect"
	fmt.Println("DetectStatus")
}
