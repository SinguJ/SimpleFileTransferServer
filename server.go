package main

import (
    "embed"
    "fmt"
    "io/fs"
    "log"
    "net/http"
    "time"
)

//go:embed static
var staticEmbedFs embed.FS

// 服务
var serv = http.NewServeMux()

func init() {
    // 注册静态页面
    staticFS, err := fs.Sub(staticEmbedFs, "static")
    if err != nil {
        log.Panicln(err)
    }
    serv.Handle("/", http.FileServer(http.FS(staticFS)))
}

func StartServer(port int) {
    server := &http.Server{
        Addr:           fmt.Sprintf(":%d", port),
        Handler:        serv,
        ReadTimeout:    30 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }
    if err := server.ListenAndServe(); err != nil {
        log.Panicln(err)
    }
}
