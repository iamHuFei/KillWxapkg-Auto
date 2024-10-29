package api

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
)

// 监控目录并通过通道返回新增的文件名
func WatchNewFiles(directory string, fileChan chan<- string) error {
	// 创建一个文件系统监听器
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	// 开始监控指定的目录
	err = watcher.Add(directory)
	if err != nil {
		return err
	}

	// 持续监听文件系统事件
	for {
		select {
		case event := <-watcher.Events:
			// 检测到新文件被创建
			if event.Op&fsnotify.Create == fsnotify.Create {
				//fmt.Printf("新文件: %s\n", event.Name)
				// 将新增的文件名发送到通道中
				fileChan <- event.Name
			}
		case err := <-watcher.Errors:
			log.Println("监控错误:", err)
		}
	}
}

func ListSubdirectories(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("无法读取目录: %w", err)
	}

	var subdirs []string

	// 遍历目录条目，筛选出子目录
	for _, entry := range entries {
		if entry.Type().IsDir() {
			subdirs = append(subdirs, entry.Name()) // 收集子目录名称
		}
	}

	return subdirs, nil
}
