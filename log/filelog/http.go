/*
 * Copyright (c) 2023. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package filelog

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	//go:embed *.html
	templateFiles embed.FS

	filesTemplate *template.Template

	upgrader = websocket.Upgrader{
		HandshakeTimeout: 0,
		ReadBufferSize:   0,
		WriteBufferSize:  0,
		WriteBufferPool:  nil,
		Subprotocols:     nil,
		Error:            nil,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: true,
	}
)

func init() {
	tps, err := template.ParseFS(templateFiles, "*")
	if err != nil {
		panic(err)
	}

	filesTemplate = tps
}
func (w *FileWriterByPeriod) ServeHTTP(prefix string) http.Handler {

	prefix = strings.Trim(prefix, "/")
	fs := fileServer{w: w}

	serveMux := http.NewServeMux()

	serveMux.HandleFunc(fmt.Sprintf("/%s/tail", prefix), fs.watch)
	serveMux.HandleFunc(fmt.Sprintf("/%s/files", prefix), fs.files)
	serveMux.HandleFunc(fmt.Sprintf("/%s/file/", prefix), fs.file(prefix))
	serveMux.HandleFunc(fmt.Sprintf("/%s/", prefix), fs.files)

	return serveMux
}

type fileServer struct {
	w *FileWriterByPeriod
}

func (f *fileServer) file(prefix string) func(w http.ResponseWriter, r *http.Request) {
	path := fmt.Sprintf("/%s/file/", prefix)
	return func(w http.ResponseWriter, r *http.Request) {
		fileController := f.w.fileController
		if fileController == nil {

			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("file log transport is not open"))
			return
		}
		fileName := strings.TrimPrefix(r.URL.Path, path)

		filePath := filepath.Join(fileController.dir, fileName)

		file, err := os.OpenFile(filePath, os.O_RDONLY, 0)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}
		grep := r.URL.Query().Get("grep")
		grepBytes := []byte(grep)

		w.Header().Set("Content-Type", "application/x-gzip")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=\"%s.gz\"", fileName))
		defer file.Close()
		buf := bufio.NewReader(file)

		flusher := w.(http.Flusher)
		lineBuf := bytes.NewBuffer(make([]byte, 0, 4096))

		zw := gzip.NewWriter(w)
		zw.Name = fileName
		zw.Comment = "apinto log file"
		defer zw.Close()
		for {
			line, b, err := buf.ReadLine()
			if len(line) > 0 {
				lineBuf.Write(line)
				if b {
					continue
				}
				lineData := lineBuf.Bytes()

				if len(grepBytes) == 0 || bytes.Contains(lineData, grepBytes) {
					zw.Write(lineData)
					zw.Write([]byte("\r\n"))
					zw.Flush()
					flusher.Flush()
				}
				lineBuf.Reset()
			}

			if err != nil {
				break
			}
		}

	}
}
func (f *fileServer) files(w http.ResponseWriter, r *http.Request) {
	fileController := f.w.fileController
	if fileController == nil {

		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("file log transport is not open"))
		return
	}
	pathPatten := filepath.Join(fileController.dir, fmt.Sprintf("%s*.log", fileController.file))
	fileMatched, err := filepath.Glob(pathPatten)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "error to list %s*.log :%s", fileController.file, err.Error())
		return
	}

	fs := make([]string, 0, len(fileMatched))
	for _, path := range fileMatched {
		fs = append(fs, filepath.Base(path))
	}

	if _, has := r.URL.Query()["show"]; has {
		w.Header().Set("Content-type", "text/html; charset=utf-8")
		filesTemplate.ExecuteTemplate(w, "files.html", fs)
		return
	}

	data, err := json.Marshal(fs)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
func (f *fileServer) watch(w http.ResponseWriter, r *http.Request) {

	if _, has := r.URL.Query()["show"]; has {
		w.Header().Set("Content-type", "text/html; charset=utf-8")

		filesTemplate.ExecuteTemplate(w, "ws.html", nil)
		return
	}

	h, err := f.w.Watch()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	defer h.Cancel()
	conn, err := upgrader.Upgrade(w, r, http.Header{})
	if err != nil {
		return
	}
	defer conn.Close()

	grep := []byte(r.URL.Query().Get("grep"))
	ctx, cancel := context.WithCancel(r.Context())
	go func() {
		defer cancel()
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}()
	for {
		select {
		case msg, ok := <-h.C:
			if !ok {
				return
			}
			if len(grep) == 0 || bytes.Contains(msg, grep) {
				err := conn.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					return
				}
			}

		case <-ctx.Done():
			return
		}
	}
}
