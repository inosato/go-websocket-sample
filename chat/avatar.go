package main

import (
	"errors"
	"fmt"
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
