package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

func concurrentZip(folders []string, numGoroutines int) {
	var wg sync.WaitGroup
	wg.Add(len(folders))

	// 创建一个channel来分发任务
	taskChan := make(chan string, numGoroutines)

	// 创建goroutines
	for i := 0; i < numGoroutines; i++ {
		go func() {
			for folder := range taskChan {
				// 创建ZIP文件名
				zipName := fmt.Sprintf("%s.zip", folder)
				err := zipFolder(folder, zipName)
				if err != nil {
					fmt.Printf("Error zipping %s: %v\n", folder, err)
				}
				wg.Done()
			}
		}()
	}

	// 分发任务
	for _, folder := range folders {
		taskChan <- folder
	}

	// 关闭channel并等待所有goroutines完成
	close(taskChan)
	wg.Wait()
}

func zipFolder(folderPath, zipPath string) error {
	// 创建一个新的zip文件
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// 创建一个zip档案的写入器
	zw := zip.NewWriter(zipFile)
	defer zw.Close()

	// 递归处理文件夹中的每个文件和子目录
	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过隐藏文件和.DS_Store文件
		if info.Name() == ".DS_Store" || (len(info.Name()) > 1 && info.Name()[0] == '.') {
			return nil
		}

		// 构建相对路径
		relPath, err := filepath.Rel(folderPath, path)
		if err != nil {
			return err
		}

		// 创建zip文件的头部
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = relPath
		header.Method = zip.Deflate

		// 如果是目录，添加斜杠以标识
		if info.IsDir() {
			if relPath != "." && relPath != ".." { // 跳过当前目录和上级目录
				header.Name += "/"
			}
		}

		// 创建zip文件的writer
		zwf, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}

		// 如果是文件，打开并复制内容到zip
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(zwf, file)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func getFirstSubDirs(inputDir string) []string {
	dirs := make([]string, 0)

	var currentDir = inputDir
	var err error
	if currentDir == "" {
		// 获取当前工作目录

		currentDir, err = os.Getwd()
		if err != nil {
			fmt.Println("Error getting current directory:", err)
			return nil
		}
	}
	// 读取当前目录下的文件和目录
	entries, err := os.ReadDir(currentDir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return dirs
	}

	// 遍历所有条目
	for _, entry := range entries {
		// 获取条目的完整路径
		path := filepath.Join(currentDir, entry.Name())

		// 检查是否是目录
		if entry.IsDir() {
			// 打印目录路径
			fmt.Println(path)
			dirs = append(dirs, path)
		}
	}
	return dirs
}
