package main

import (
	"fmt"
	"github.com/kardianos/service"
	"os"
	"path/filepath"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func (p *program) run() {
	// Do work here
	main()
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func serviceInstall(isInstall bool) {
	cwd, _ := os.Getwd()
	wd, _ := filepath.Abs(cwd)

	svcConfig := &service.Config{
		Name:             "sealchat-backend",
		DisplayName:      "sealchat-backend",
		Description:      "sealchat-backend",
		WorkingDirectory: wd,
	}

	prg := &program{}
	fmt.Println("正在试图访问系统服务 ...")
	s, err := service.New(prg, svcConfig)

	if isInstall {
		fmt.Println("正在安装系统服务，安装完成后，SealDice将自动随系统启动")
		if err != nil {
			fmt.Printf("安装失败: %s\n", err.Error())
		}
		_, err = s.Logger(nil)
		if err != nil {
			fmt.Printf("安装失败: %s\n", err.Error())
			fmt.Println(err)
		}
		err = s.Install()
		if err != nil {
			fmt.Printf("安装失败: %s\n", err.Error())
			return
		}

		fmt.Println("安装完成，正在启动……")
		s.Start()
	} else {
		fmt.Println("正在卸载系统服务……")
		s.Stop()
		s.Uninstall()
		fmt.Println("系统服务已删除")
	}
}
