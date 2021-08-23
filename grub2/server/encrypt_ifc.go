package server

import (
	"fmt"
	"grub-encrypt-dbus/grub2/utils"

	"github.com/godbus/dbus"
	polkit "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.policykit1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/log"
)

// dbus服务信息
const (
	dbusName = "com.deepin.daemon.GrubEncryption"  // @Name: 	名称
	dbusPath = "/com/deepin/daemon/GrubEncryption" // @Path:	地址
	dbusIFC  = "com.deepin.daemon.GrubEncryption"  // @IFC:	接口名
)

var logger *log.Logger

func SetLogger(v *log.Logger) {
	logger = v
}

// 获取 dbus对象 ifc名称
func (g *GrubEncrypt) GetInterfaceName() string {
	return dbusIFC
}

//创建菜单加密账户，创建
func (g *GrubEncrypt) AddAccount(sender dbus.Sender, usr, passwd string) *dbus.Error {
	if utils.IsInSlice(usr, g.OnlineUser) || utils.IsInSlice(usr, g.OfflineUser) {
		logger.Warning("Existing user: " + usr)
		return dbusutil.ToError(fmt.Errorf("Existing user: %s\n", usr))
	}
	chiperPasswd, err := utils.GrubPBKDF2Crypto(usr, passwd)
	if err != nil {
		fmt.Println(err)
		return dbusutil.ToError(err)
	}
	err = utils.Add.WriteConfig(usr, chiperPasswd)
	if err != nil {
		fmt.Println(err)
		return dbusutil.ToError(err)
	}
	g.DeatectAccount()
	logger.Info("Add user: " + usr)
	//更新 update-gurb
	return nil
}

func (g *GrubEncrypt) DeleteAccount(usr string) *dbus.Error {
	err := utils.Delete.WriteConfig(usr)
	if err != nil {
		fmt.Println(err)
		return dbusutil.ToError(err)
	}
	g.DeatectAccount()
	logger.Info("Delete user: " + usr)
	//更新 update-gurb

	return nil
}

//失能所有用户
func (g *GrubEncrypt) DisableAuthentication(usr string) *dbus.Error {
	err := utils.Disable.WriteConfig(usr)
	if err != nil {
		fmt.Println(err)
		return dbusutil.ToError(err)
	}
	g.Status = "disable"
	g.DeatectAccount()
	logger.Info("disable grub menu encryption")
	return nil
}

//使能所有用户
func (g *GrubEncrypt) EnableAuthentication(sender dbus.Sender, usr string) *dbus.Error {
	// fmt.Println(sender)
	_, ok := g.checkAuth(sender, "com.deepin.daemon.GrubEncryption")
	if ok != nil {
		fmt.Println(ok.Error())
	}

	err := utils.Enable.WriteConfig(usr)
	if err != nil {
		fmt.Println(err)
		return dbusutil.ToError(err)
	}
	g.Status = "enable"
	g.DeatectAccount()
	logger.Info("enable grub menu encryption")
	return nil
}

//检测当前账户
func (g *GrubEncrypt) DeatectAccount() {
	onlineUserList, offlineUserList, err := utils.DoDetect()
	if err != nil {
		fmt.Println(err)
	}
	g.OnlineUser = onlineUserList
	g.OfflineUser = offlineUserList
	if g.OnlineUser != nil {
		g.Status = "enable"
	} else {
		g.Status = "disable"
	}
}

//鉴权
func (g *GrubEncrypt) checkAuth(sender dbus.Sender, actionId string) (bool, error) {
	systemBus, err := dbus.SystemBus()
	if err != nil {
		return false, err
	}
	authority := polkit.NewAuthority(systemBus)
	subject := polkit.MakeSubject(polkit.SubjectKindSystemBusName)
	subject.SetDetail("name", string(sender))
	result, err := authority.CheckAuthorization(0, subject, actionId, nil,
		polkit.CheckAuthorizationFlagsAllowUserInteraction, "")
	if err != nil {
		return false, err
	}

	return result.IsAuthorized, nil
}
