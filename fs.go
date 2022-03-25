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

/** * 结构体 *  **/

// TargetType 目标类型
type TargetType int

const (
    // TypeFile 文件
    TypeFile TargetType = 0
    // TypeDir 目录（文件夹）
    TypeDir = 1
)

// Metadata 元数据
type Metadata interface {
    // GetName 获取文件（夹）名称
    GetName() string
    // GetPath 获取文件（夹）路径
    GetPath() string
    // GetType 获取目标类型
    GetType() TargetType
    // AsFile 作为文件
    AsFile() *File
    // AsDir 作为目录
    AsDir() *Directory
}

// MetadataImpl 元数据
type MetadataImpl struct {
    // 文件（夹）名称
    Name string `json:"name"`
    // 文件（夹）路径
    Path string `json:"path"`

    // 目标类型
    Type TargetType `json:"type"`
}

func (m *MetadataImpl) GetName() string {
    return m.Name
}

func (m *MetadataImpl) GetPath() string {
    return m.Path
}

func (m *MetadataImpl) GetType() TargetType {
    return m.Type
}

func (m *MetadataImpl) AsFile() *File {
    panic("implement me")
}

func (m *MetadataImpl) AsDir() *Directory {
    panic("implement me")
}

// File 文件
type File struct {
    MetadataImpl

    // 文件大小
    Size int64 `json:"size"`
    // 文件类型
    Filetype string `json:"filetype"`
}

func (f *File) AsFile() *File {
    return f
}

func (f *File) AsDir() *Directory {
    panic("wrong call")
}

// Directory 目录（文件夹）
type Directory struct {
    MetadataImpl

    // 子文件列表
    Files []Metadata `json:"files"`
}

func (d *Directory) AsFile() *File {
    panic("wrong call")
}

func (d *Directory) AsDir() *Directory {
    return d
}
