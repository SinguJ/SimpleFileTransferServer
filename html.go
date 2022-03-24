package main

import (
	"embed"
	_ "embed"
	"log"
	"net/http"
	"path/filepath"
	templateLib "text/template"
)

//go:embed html
var tmpl embed.FS

var template *templateLib.Template

func init() {
	var err error
	template, err = templateLib.ParseFS(tmpl, "html/*.html")
	if err != nil {
		log.Panicln(err)
	}
}

type FileInfo struct {
	Name  string
	Path  string
	IsDir bool
}

func viewFileList(writer http.ResponseWriter, dir string) error {
	// 获取文件列表
	files, err := getFiles(dir)
	if err != nil {
		return err
	}
	// 构建渲染数据
	fileInfos := make([]*FileInfo, len(files))
	for index, entry := range files {
		filename := entry.Name()
		fileInfos[index] = &FileInfo{
			Name:  filename,
			Path:  filepath.Join(dir, filename),
			IsDir: entry.IsDir(),
		}
	}
	// 渲染
	return template.ExecuteTemplate(writer, "file_list.html", fileInfos)
}
