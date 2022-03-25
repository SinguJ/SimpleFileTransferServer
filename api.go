package main

import (
    "errors"
    "fmt"
    "log"
    "net/http"
    "net/url"
    "os"
    "path"
    "strings"
)

func init() {
    // 注册下载接口
    serv.HandleFunc("/api/download/", apiHandle(func(writer http.ResponseWriter, request *http.Request) {
        // if strings.HasPrefix(request.RequestURI, "/download") {
        //     download(writer, request.RequestURI[9:])
        // }
        // show(writer, request.RequestURI)
        p, err := readTargetPath(request, "/api/download/")
        if err != nil {
            panic(err)
        }
        downloadOrView(writer, p)
    }))
}

// ErrServiceComplete 服务完成
var ErrServiceComplete = errors.New("服务完成")

// ErrNotFound 找不到指定的文件或文件夹
var ErrNotFound = errors.New("找不到指定的文件或文件夹")

// 读取请求的文件路径
func readTargetPath(r *http.Request, prefix string) (targetPath string, err error) {
    // 解析 URI
    p, err := url.QueryUnescape(r.URL.Path)
    if err != nil {
        return "", err
    }
    // 去除接口前缀
    if prefix != "" {
        p = strings.TrimPrefix(p, prefix)
    }

    if p == "" {
        targetPath = "/"
        return
    }
    if p[0] != '/' {
        p = "/" + p
    }
    targetPath = path.Clean(p)
    // path.Clean removes trailing slash except for root;
    // put the trailing slash back if necessary.
    if p[len(p)-1] == '/' && targetPath != "/" {
        // Fast path for common case of p being the string we want:
        if len(p) == len(targetPath)+1 && strings.HasPrefix(p, targetPath) {
            targetPath = p
        } else {
            targetPath += "/"
        }
    }
    return
}

// 下载
func download(writer http.ResponseWriter, p string) {
    var err error

    // 获取路径的状态
    stat, err := os.Stat(p)
    if err != nil {
        if os.IsNotExist(err) {
            panic(ErrNotFound)
        }
        panic(err)
    }
    // 检查是否是目录
    if stat.IsDir() {
        // TODO: 压缩下载
        panic("暂不支持下载文件夹")
    } else {
        fileWriteTo(writer, p)
    }
}

// 展示文件或目录
func show(writer http.ResponseWriter, p string) {
    var err error

    // 获取路径的状态
    stat, err := os.Stat(p)
    if err != nil {
        if os.IsNotExist(err) {
            panic(ErrNotFound)
        }
        panic(err)
    }
    // 检查是否是目录
    if stat.IsDir() {
        // 展示文件列表
        _, err := getFiles(p)
        if err != nil {
            return
        }
    } else {
        // TODO: 展示文件摘要
        panic("暂不支持查看文件信息")
    }
}

// 临时解决方法
func downloadOrView(writer http.ResponseWriter, p string) {
    var err error
    _path := getAbsPath(p)

    // 获取路径的状态
    stat, err := os.Stat(_path)
    if err != nil {
        if os.IsNotExist(err) {
            panic(ErrNotFound)
        }
        panic(err)
    }
    // 检查是否是目录
    if stat.IsDir() {
        // 不允许下载文件夹
        panic("暂不支持下载文件夹")
    } else {
        // TODO: 展示文件摘要
        fileWriteTo(writer, p)
    }
}

// 接口处理装饰器函数
func apiHandle(handle func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            err := recover()
            if err != nil {
                switch err {
                case ErrServiceComplete:
                    return
                default:
                    _, _err := fmt.Fprintln(w, "程序错误：", err)
                    if _err != nil {
                        log.Println(_err)
                    }
                }
            }
        }()
        handle(w, r)
    }
}
