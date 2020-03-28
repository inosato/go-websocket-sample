package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/inosato/go-websocket-sample/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/objx"
	"github.com/stretchr/signature"
)

var avatars Avatar = UserFileSystemAvatar

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "address of the application")
	flag.Parse()

	gomniauth.SetSecurityKey(signature.RandomKey(16))
	gomniauth.WithProviders(
		github.New(os.Getenv("GITHUB_ID"), os.Getenv("GITHUB_SECRET"), os.Getenv("GITHUB_REDIRECT")),
	)

	// r := newRoom(UseAuthAvatar)
	// r := newRoom(UseGravatarAvatar)
	r := newRoom(UserFileSystemAvatar)
	r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.Handle("/room", r)
	http.HandleFunc("/auth/", loginHandler)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.HandleFunc("/uploader", UploaderHandler)
	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars"))))
	log.Println("Start WebSocket")
	go r.run()
	log.Println("Start WebServer Port:", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
