package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Ackites/KillWxapkg/api"
	"github.com/Ackites/KillWxapkg/cmd"
	hook2 "github.com/Ackites/KillWxapkg/internal/hook"
	"github.com/Ackites/KillWxapkg/internal/pack"
	"log"
	"path/filepath"
	"time"
)

func test() {
	config, err := api.LoadConfig("config/config.yaml")
	if err != nil {
		return
	}
	// 定义要监控的目录
	directory := config.WxAppletPath // 请确保该目录存在

	// 创建一个通道用于接收文件名
	fileChan := make(chan string)
	fmt.Printf("监视目录：%s\n", directory)
	// 启动一个goroutine进行文件监控
	go func() {
		err := api.WatchNewFiles(directory, fileChan)
		if err != nil {
			log.Fatalf("无法启动文件监控: %v", err)
		}
	}()

	// 主线程持续接收新增文件名并打印
	for {
		newFile := <-fileChan
		fmt.Printf("接收新增文件: %s\n", newFile)
		// 等待目录加载时间
		time.Sleep(10 * time.Second)
		fmt.Println("等待加载目录.")
		subdirs, err := api.ListSubdirectories(newFile)
		if err != nil {
			log.Fatalf("错误: %v", err)
		}

		// 发送请求获取appid信息，不一定成功，需要判断
		body, err := api.SendRequest(filepath.Base(newFile))
		if err != nil {
			fmt.Println(err)
			return
		}
		var responseData api.ResponseData
		if err := json.Unmarshal(body, &responseData); err != nil {
			fmt.Println("解析JSON失败:", err)
			return
		}

		appID := filepath.Base(newFile)
		input := newFile + "\\" + subdirs[0] + "\\__APP__.wxapkg"
		outputDir := "output\\"
		if responseData.Data.Name != "" {
			outputDir += responseData.Data.Name
		} else {
			outputDir += appID
		}
		fmt.Println("准备输出目录:", outputDir)
		fileExt := ""
		restoreDir := false
		pretty := false
		noClean := true
		save := false
		sensitive := true
		cmd.Execute(appID, input, outputDir, fileExt, restoreDir, pretty, noClean, save, sensitive)
		fmt.Println(">>> 完成逆向编译<-" + responseData.Data.Name)
		fmt.Println(">>> 等待下一任务->")
	}
	//cmd.Execute(appID, input, outputDir, fileExt, restoreDir, pretty, noClean, save, sensitive)

}

var (
	auto       bool
	appID      string
	input      string
	outputDir  string
	fileExt    string
	restoreDir bool
	pretty     bool
	noClean    bool
	hook       bool
	save       bool
	repack     string
	watch      bool
	sensitive  bool
)

func init() {
	flag.BoolVar(&auto, "auto", false, "是否目录监控自动反编译，点击即是反编译")
	flag.StringVar(&appID, "id", "", "微信小程序的AppID")
	flag.StringVar(&input, "in", "", "输入文件路径（多个文件用逗号分隔）或输入目录路径")
	flag.StringVar(&outputDir, "out", "", "输出目录路径（如果未指定，则默认保存到输入目录下以AppID命名的文件夹）")
	flag.StringVar(&fileExt, "ext", ".wxapkg", "处理的文件后缀")
	flag.BoolVar(&restoreDir, "restore", false, "是否还原工程目录结构")
	flag.BoolVar(&pretty, "pretty", false, "是否美化输出")
	flag.BoolVar(&noClean, "noClean", false, "是否清理中间文件")
	flag.BoolVar(&hook, "hook", false, "是否开启动态调试")
	flag.BoolVar(&save, "save", false, "是否保存解密后的文件")
	flag.StringVar(&repack, "repack", "", "重新打包wxapkg文件")
	flag.BoolVar(&watch, "watch", false, "是否监听将要打包的文件夹，并自动打包")
	flag.BoolVar(&sensitive, "sensitive", false, "是否获取敏感数据")
}

func main() {
	// 解析命令行参数
	flag.Parse()

	banner := `

 ____  __.__.______   _________          __          
|    |/ _|__|  \   \ /   /  _  \  __ ___/  |_  ____  
|      < |  |  |\   Y   /  /_\  \|  |  \   __\/  _ \ 
|    |  \|  |  |_\     /    |    \  |  /|  | (  <_> )
|____|__ \__|____/\___/\____|__  /____/ |__|  \____/ 
        \/                     \/
                                                    
             Wxapkg Decompiler Tool v1.0.0
    `
	fmt.Println(banner)
	// 自动处理
	if auto {
		test()
		return
	}
	// 动态调试
	if hook {
		hook2.Hook()
		return
	}

	// 重新打包
	if repack != "" {
		pack.Repack(repack, watch, outputDir)
		return
	}

	// 参数检查
	if appID == "" || input == "" {
		fmt.Println("使用方法: program -id=<AppID> -in=<输入文件1,输入文件2> 或 -in=<输入目录> -out=<输出目录> [-ext=<文件后缀>] [-restore] [-pretty] [-noClean] [-hook] [-save] [-repack=<输入目录>] [-watch] [-sensitive]")
		flag.PrintDefaults()
		fmt.Println()
		return
	}

	// 执行命令
	cmd.Execute(appID, input, outputDir, fileExt, restoreDir, pretty, noClean, save, sensitive)
	//cmd.Execute("wxd4185d00bf7e08ac", "D:\\appdata\\tencent\\WeChat\\documents\\WeChat Files\\Applet\\wxd4185d00bf7e08ac\\551\\__APP__.wxapkg", "output", fileExt, restoreDir, pretty, noClean, save, sensitive)
}
