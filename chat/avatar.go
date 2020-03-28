package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

var ErrNoAvatar = errors.New("chat: cannot get avatar")

type Avatar interface {
	GetAvatarURL(u ChatUser) (string, error)
}

type TryAvatars []Avatar

func (a TryAvatars) GetAvatarURL(u ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(u); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatar
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (_ AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	if url := u.AvatarURL(); len(url) != 0 {
		return url, nil
	}
	return "", ErrNoAvatar
}

type GravatarAvatar struct{}

var UseGravatarAvatar GravatarAvatar

func (_ GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return fmt.Sprintf("//www.gravatar.com/avatar/%s", u.UniqueID()), nil
}

type FileSystemAvatar struct{}

var UserFileSystemAvatar FileSystemAvatar

func (_ FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	if files, err := ioutil.ReadDir("avatars"); err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if match, _ := filepath.Match(u.UniqueID()+"*", file.Name()); match {
				return fmt.Sprintf("/avatars/%s", file.Name()), nil
			}
		}
	}

	return "", ErrNoAvatar
}
