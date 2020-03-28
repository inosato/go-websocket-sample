package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

var ErrNoAvatar = errors.New("chat: cannot get avatar")

type Avatar interface {
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (_ AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"].(string); ok {
		return url, nil
	}
	return "", ErrNoAvatar
}

type GravatarAvatar struct{}

var UseGravatarAvatar GravatarAvatar

func (_ GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if userID, ok := c.userData["userid"]; ok {
		if userIDStr, ok := userID.(string); ok {
			return fmt.Sprintf("//www.gravatar.com/avatar/%s", userIDStr), nil
		}
	}
	return "", ErrNoAvatar
}

type FileSystemAvatar struct{}

var UserFileSystemAvatar FileSystemAvatar

func (_ FileSystemAvatar) GetAvatarURL(c *client) (string, error) {
	if userID, ok := c.userData["userid"]; ok {
		if userIDStr, ok := userID.(string); ok {
			if files, err := ioutil.ReadDir("avatars"); err == nil {
				for _, file := range files {
					if file.IsDir() {
						continue
					}
					if match, _ := filepath.Match(userIDStr+"*", file.Name()); match {
						return fmt.Sprintf("/avatars/%s", file.Name()), nil
					}
				}
			}

		}
	}
	return "", ErrNoAvatar
}
