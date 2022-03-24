package main

import (
    "embed"
    "errors"
    "fmt"
    "io/fs"
    "log"
    "net/http"
    "net/url"
    "os"
    "time"
)

// ErrServiceOver 服务结束
var ErrServiceOver = errors.New("服务结束")

// 404
func notFound(writer http.ResponseWriter, path string) {
    _, err := fmt.Fprintln(writer, "找不到指定的文件或目录：", path[1:])
    if err != nil {
        log.Println(err)
    }
    panic(ErrServiceOver)
}

// 500
func serviceError(err interface{}) {
    panic(err)
}

// 下载
func download(writer http.ResponseWriter, path string) {
    var err error

    // 获取路径的状态
    stat, err := os.Stat(path)
    if err != nil {
        if os.IsNotExist(err) {
            notFound(writer, path)
        }
        serviceError(err)
    }
    // 检查是否是目录
    if stat.IsDir() {
        // TODO: 压缩下载
        serviceError("暂不支持下载文件夹")
    } else {
        fileWriteTo(writer, path)
    }
}

// 展示文件或目录
func show(writer http.ResponseWriter, path string) {
    var err error

    // 获取路径的状态
    stat, err := os.Stat(path)
    if err != nil {
        if os.IsNotExist(err) {
            notFound(writer, path)
        }
        serviceError(err)
    }
    // 检查是否是目录
    if stat.IsDir() {
        // 展示文件列表
        _, err := getFiles(path)
        if err != nil {
            return
        }
    } else {
        // TODO: 展示文件摘要
        serviceError("暂不支持查看文件信息")
    }
}

// 临时解决方法
func downloadOrView(writer http.ResponseWriter, path string) {
    var err error
    _path := getAbsPath(path)

    // 获取路径的状态
    stat, err := os.Stat(_path)
    if err != nil {
        if os.IsNotExist(err) {
            notFound(writer, path)
        }
        serviceError(err)
    }
    // 检查是否是目录
    if stat.IsDir() {
        // 展示文件列表
        err = viewFileList(writer, path)
        if err == nil {
            err = ErrServiceOver
        }
        serviceError(err)
    } else {
        // TODO: 展示文件摘要
        fileWriteTo(writer, path)
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
    path, err := url.QueryUnescape(request.RequestURI)
    if err != nil {
        serviceError(err)
    }
    if path == "" {
        path = "/"
    }
    downloadOrView(writer, path)
}

//go:embed static
var staticEmbedFs embed.FS

// 注册 HTTP 处理过程
func registerHttpHandles() {
    // 注册静态页面
    staticFS, err := fs.Sub(staticEmbedFs, "static")
    if err != nil {
        log.Panicln(err)
    }
    http.Handle("/", http.FileServer(http.FS(staticFS)))
}

func StartServer(port int) {
    registerHttpHandles()
    server := &http.Server{
        Addr:           fmt.Sprintf(":%d", port),
        Handler:        http.DefaultServeMux,
        ReadTimeout:    30 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }
    if err := server.ListenAndServe(); err != nil {
        log.Panicln(err)
    }
}
