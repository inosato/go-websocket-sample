package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
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

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	client := new(client)
	client.userData = map[string]interface{}{
		"userid": "686f6765686f6765",
	}

	url, err := gravatarAvatar.GetAvatarURL(client)
	assert.NoError(t, err)

	assert.Equal(t, "//www.gravatar.com/avatar/686f6765686f6765", url)
}

func TestFileSystemAvatar(t *testing.T) {
	filename := filepath.Join("avatars", "abc.png")
	ioutil.WriteFile(filename, []byte{}, 0777)
	defer func() { os.Remove(filename) }()

	var fileSystemAvatar FileSystemAvatar
	client := new(client)
	client.userData = map[string]interface{}{
		"userid": "xyz",
	}
	url, err := fileSystemAvatar.GetAvatarURL(client)
	assert.NoError(t, err)
	assert.Equal(t, "/avatars/xyz.png", url)
}
