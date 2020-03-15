package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar

	client := new(client)
	_, err := authAvatar.GetAvatarURL(client)
	assert.Error(t, err, "chat: cannot get avatar")

	testURL := "http://url-to-avatar/"
	client.userData = map[string]interface{}{
		"avatar_url": testURL,
	}

	url, err := authAvatar.GetAvatarURL(client)
	assert.Equal(t, testURL, url)

}
