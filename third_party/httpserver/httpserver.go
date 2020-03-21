package httpserver

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

func Run(addr string) (err error) {
	_, filename, _, _ := runtime.Caller(0)
	static := path.Join(filepath.Dir(filename), "static")
	fileSrv := http.FileServer(http.Dir(static))
	http.Handle("/static/", http.StripPrefix("/static/", fileSrv))
	http.HandleFunc("/", serveTemplate)

	if err := http.ListenAndServe(addr, nil); err != nil {
		return err
	}
	return nil
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	_, filename, _, _ := runtime.Caller(0)
	tmpDir := filepath.Dir(filename)
	tplFile := r.URL.Path
	if tplFile == "/" {
		tplFile = "/index.html"
	}
	lp := filepath.Join(tmpDir, "templates", "_layout.html")
	fp := filepath.Join(tmpDir, "templates", tplFile)
	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

func Addr() string {
	port, err := Port()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("0.0.0.0:%d", port)
}

func Port() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
