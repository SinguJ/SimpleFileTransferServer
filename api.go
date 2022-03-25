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
    serv.Handle("/api/download/", &Handler{})
}

// ErrServiceOver 服务结束
var ErrServiceOver = errors.New("服务结束")

// 404
func notFound(writer http.ResponseWriter, p string) {
    _, err := fmt.Fprintln(writer, "找不到指定的文件或目录：", p[1:])
    if err != nil {
        log.Println(err)
    }
    panic(ErrServiceOver)
}

// 500
func serviceError(err interface{}) {
    panic(err)
}

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
            notFound(writer, p)
        }
        serviceError(err)
    }
    // 检查是否是目录
    if stat.IsDir() {
        // TODO: 压缩下载
        serviceError("暂不支持下载文件夹")
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
            notFound(writer, p)
        }
        serviceError(err)
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
        serviceError("暂不支持查看文件信息")
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
            notFound(writer, p)
        }
        serviceError(err)
    }
    // 检查是否是目录
    if stat.IsDir() {
        // 展示文件列表
        // err = viewFileList(writer, p)
        // if err == nil {
        //     err = ErrServiceOver
        // }
        // serviceError(err)
    } else {
        // TODO: 展示文件摘要
        fileWriteTo(writer, p)
    }
}

type Handler struct{}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
    defer func() {
        err := recover()
        if err != nil && err != ErrServiceOver {
            _, _err := fmt.Fprintln(writer, "程序错误：", err)
            if _err != nil {
                log.Println(_err)
            }
        }
    }()

    // if strings.HasPrefix(request.RequestURI, "/download") {
    //     download(writer, request.RequestURI[9:])
    // }
    // show(writer, request.RequestURI)
    p, err := readTargetPath(request, "/api/download/")
    if err != nil {
        serviceError(err)
    }
    downloadOrView(writer, p)
}
