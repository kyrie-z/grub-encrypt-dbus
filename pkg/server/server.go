package server

import (
	"fmt"
	"grub-dbus/pkg/utils"

	"pkg.deepin.io/lib/dbusutil"
)

/*
	提供dbus服务信息的相关定义。
*/

// dbus服务信息
const (
	dbusName = "com.deepin.daemon.GrubEncryption"  // @Name: 	名称
	dbusPath = "/com/deepin/daemon/GrubEncryption" // @Path:	地址
	dbusIFC  = "com.deepin.daemon.GrubEncryption"  // @IFC:	接口名
)

type GrubEncrypt struct {
	methods *struct {
		Setpassword   func() `in:"user,password" out:"user"`
		Unsetpassword func() `in:"user" out:"user"`
	}
	Status string
}

type Service struct {
	conn        *dbusutil.Service
	GrubEncrypt *GrubEncrypt
}

var dbusSrv *Service

// 新建对象
func newService() (*Service, error) {
	srv, err := dbusutil.NewSystemService()
	if err != nil {
		return nil, fmt.Errorf("new system service is error:%s\n", err)
	}
	grubEncrypt := &GrubEncrypt{}
	return &Service{conn: srv, GrubEncrypt: grubEncrypt}, nil
}

// 获取初始化的 dbus 对象,不存在就新建
func GetService() *Service {
	if dbusSrv != nil {
		return dbusSrv
	}
	var err error
	dbusSrv, err = newService()
	if err != nil {
		panic(err)
	}
	return dbusSrv
}

// 获取 dbus对象 ifc名称
func (g *GrubEncrypt) GetInterfaceName() string {
	return dbusIFC
}

// 外部调用
func (s *Service) Init() error {
	err := s.conn.Export(dbusPath, s.GrubEncrypt)
	if err != nil {
		return err
	}
	return s.conn.RequestName(dbusName)
}

func (s *Service) Loop() {
	//变量初始化操作
	s.GrubEncrypt.detectStatus()
	utils.InitGrubPBKDF2()

	//wait
	s.conn.Wait()
}
