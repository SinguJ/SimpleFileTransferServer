package main

import (
    "encoding/json"
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
    serv.HandleFunc("/api/download/", fileHandle("/api/download/", download))
    // 注册文件信息接口
    serv.HandleFunc("/api/metadata/", fileHandle("/api/metadata/", metadata))
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
func download(w http.ResponseWriter, _ *http.Request, p string, stat os.FileInfo) {
    // 检查是否是目录
    if stat.IsDir() {
        panic("暂不支持下载文件夹")
    } else {
        // 写出文件
        fileWriteTo(w, p)
    }
}

// 展示文件或文件夹信息
func metadata(w http.ResponseWriter, _ *http.Request, p string, _ os.FileInfo) {
    // 读取元数据
    metadata, err := getMetadata(p, 1)
    if err != nil {
        panic(err)
    }
    // 响应 JSON 数据
    success(w, metadata)
}

// 文件处理装饰器函数
func fileHandle(uriPrefix string, handle func(http.ResponseWriter, *http.Request, string, os.FileInfo)) func(http.ResponseWriter, *http.Request) {
    return apiHandle(func(w http.ResponseWriter, r *http.Request) {
        // 读取请求的目标路径
        p, err := readTargetPath(r, uriPrefix)
        if err != nil {
            panic(err)
        }
        // 获取目标文件在文件系统中的真实路径
        realPath := getAbsPath(p)
        // 检查文件状态
        stat, err := os.Stat(realPath)
        if err != nil {
            if os.IsNotExist(err) {
                panic(ErrNotFound)
            }
            panic(err)
        }
        handle(w, r, p, stat)
    })
}

// 接口处理装饰器函数
func apiHandle(handle func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            err := recover()
            if err != nil {
                switch err {
                case ErrServiceComplete:
                    // 正常结束的服务
                    success(w, nil)
                case ErrNotFound:
                    // 找不到文件的错误
                    fail(w, 404, ErrNotFound)
                default:
                    // 其他错误
                    log.Println(err)
                    fail(w, 500, err)
                }
            }
        }()
        handle(w, r)
    }
}

// Response 接口响应体
type Response struct {
    // 状态码
    Code int `json:"code"`
    // 消息
    Message string `json:"message"`
    // 结果
    Result interface{} `json:"result"`
}

// 成功
func success(w http.ResponseWriter, result interface{}) {
    responseJsonData(w, 0, "", result)
}

// 失败
func fail(w http.ResponseWriter, code int, msg interface{}) {
    responseJsonData(w, code, fmt.Sprint(msg), nil)
}

// 响应 JSON 数据
func responseJsonData(w http.ResponseWriter, code int, message string, result interface{}) {
    resp := &Response{
        Code:    code,
        Message: message,
        Result:  result,
    }
    data, err := json.Marshal(resp)
    if err != nil {
        log.Println(message)
        log.Println(err)
        data = []byte(`{"code": -1, "message": "服务异常", "result": null}`)
    }
    _, err = w.Write(data)
    if err != nil {
        log.Println(message)
        log.Println(err)
    }
}
