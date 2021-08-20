package server

import (
	"fmt"
	"grub-encrypt-dbus/grub2/utils"

	"pkg.deepin.io/lib/dbusutil"
)

/*
	提供dbus服务信息的相关定义。
*/

type GrubEncrypt struct {
	methods *struct {
		AddAccount            func() `in:"user,password" out:"user"`
		DeleteAccount         func() `in:"user"`
		DisableAuthentication func() `in:"user"`
		EnableAuthentication  func() `in:"user"`
	}
	Status      string
	OnlineUser  []string
	OfflineUser []string
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

func (s *Service) Init() error {
	err := s.conn.Export(dbusPath, s.GrubEncrypt)
	if err != nil {
		return err
	}
	return s.conn.RequestName(dbusName)
}

func (s *Service) Loop() {
	//变量初始化操作
	s.GrubEncrypt.DeatectAccount()
	utils.InitGrubPBKDF2()

	//wait
	s.conn.Wait()
}
