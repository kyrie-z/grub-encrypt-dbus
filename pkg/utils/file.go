package utils

import (
	"bufio"
	"errors"
	"os"
	"strings"
	"sync"
)

/*
	文件操作
*/

const (
	// const MenuCryptoCfg = "../../etc/grub.d/42_uos_menu_crypto"
	MenuCryptoCfg            = "./42_uos_menu_crypto"
	Add           ActionType = iota + 1
	Delete
	Disable
	Enable
)

var fileLock sync.Mutex

type ActionType int

//写入文件
func (action ActionType) WriteConfig(usr string, chiperPasswd ...string) error {
	fileLock.Lock()
	switch action {
	case Add:
		{
			file, err := os.OpenFile(MenuCryptoCfg, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
			if err != nil {
				return err
			}
			defer file.Close()
			write := bufio.NewWriter(file)
			write.WriteString("set superusers=\"" + usr + "\"\n")
			write.WriteString(chiperPasswd[0] + "\n")
			write.Flush()

		}
	case Delete, Disable, Enable:
		err := wirteType(usr, action)
		if err != nil {
			return err
		}
	default:
		return errors.New("undefined operation")

	}
	fileLock.Unlock()
	return nil
}

func wirteType(usr string, action ActionType) error {
	lines, err := readline(MenuCryptoCfg)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(MenuCryptoCfg, os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	for _, line := range lines {
		if strings.Contains(line, "set superusers=\""+usr+"\"") ||
			strings.Contains(line, "password_pbkdf2 "+usr+" grub.pbkdf2") {
			switch action {
			case Delete:
				continue //跳过

			case Disable:
				write.WriteString("#" + line) //添加注释

			case Enable:
				write.WriteString(strings.Trim(line, "#")) //删除备注
			}
		} else {
			write.WriteString(line)
		}
	}
	write.Flush()
	return nil
}

func readline(file string) ([]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	var lines []string
	for {
		line, _ := r.ReadString('\n')
		if line == "" {
			break
		}
		lines = append(lines, line)
	}
	return lines, nil
}
