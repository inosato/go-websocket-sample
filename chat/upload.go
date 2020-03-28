package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

func UploaderHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("UploaderHandler start")
	userID := r.FormValue("userid")

	file, header, err := r.FormFile("avatarFile")
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	filename := filepath.Join("avatars", userID+filepath.Ext(header.Filename))
	log.Print(filename)
	if err = ioutil.WriteFile(filename, data, 0777); err != nil {
		log.Print(err)
		io.WriteString(w, err.Error())
		return
	}
	io.WriteString(w, "success!")
}
