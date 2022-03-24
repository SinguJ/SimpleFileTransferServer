package main

import (
    "errors"
    "io"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
)

var rootPath string

func init() {
    var err error
    rootPath, err = os.Getwd()
    if err != nil {
        log.Panicln(err)
    }
}

var ErrNotDir = errors.New("指定的目录不是文件夹")

func getAbsPath(path string) string {
    return filepath.Join(rootPath, path)
}

// 判断文件是否是文件夹
func isFile(path string) bool {
    var err error
    // 获取路径的状态
    stat, err := os.Stat(path)
    if err != nil {
        return false
    }
    // 检查是否是目录
    return !stat.IsDir()
}

// 判断文件是否是文件夹
func isDir(path string) bool {
    var err error
    // 获取路径的状态
    stat, err := os.Stat(path)
    if err != nil {
        return false
    }
    // 检查是否是目录
    return stat.IsDir()
}

// 获取指定路径下的所有文件
func getFiles(path string) ([]os.DirEntry, error) {
    // 计算完整路径
    _path := filepath.Join(rootPath, path)

    // 检查文件是否是文件夹
    if !isDir(_path) {
        return nil, ErrNotDir
    }

    // 读取目录
    return os.ReadDir(_path)
}

func fileWriteTo(writer http.ResponseWriter, path string) {
    // 计算完整路径
    _path := filepath.Join(rootPath, path)

    _, filename := filepath.Split(_path)

    file, _ := os.Open(_path)
    defer func(file *os.File) {
        err := file.Close()
        if err != nil {
            log.Println(err)
        }
    }(file)

    fileHeader := make([]byte, 512)
    _, err := file.Read(fileHeader)
    if err != nil {
        log.Panicln(err)
    }

    fileStat, _ := file.Stat()

    writer.Header().Set("Content-Disposition", "attachment; filename="+filename)
    writer.Header().Set("Content-Type", http.DetectContentType(fileHeader))
    writer.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))

    _, err = file.Seek(0, 0)
    if err != nil {
        log.Panicln(err)
    }
    _, err = io.Copy(writer, file)
    if err != nil {
        log.Panicln(err)
    }
}
