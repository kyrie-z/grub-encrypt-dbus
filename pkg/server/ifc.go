package server

import (
	"fmt"
	"grub-dbus/pkg/utils"

	"github.com/godbus/dbus"
	"pkg.deepin.io/lib/dbusutil"
)

func (g *GrubEncrypt) Setpassword(usr, passwd string) *dbus.Error {
	chiperPasswd, err := utils.GrubPBKDF2Crypto(usr, passwd)
	if err != nil {
		fmt.Println(err)
		return dbusutil.ToError(err)
	}
	fmt.Println(chiperPasswd)
	//写入配置文件（检测是否存在帐号密码，）
	g.Status = "enable"
	return nil
}

func (g *GrubEncrypt) Unsetpassword(usr string) *dbus.Error {
	fmt.Println("delete: ", usr)
	g.Status = "disable"
	return nil
}

func (g *GrubEncrypt) detectStatus() {
	g.Status = "detect"
	fmt.Println("DetectStatus")
}
