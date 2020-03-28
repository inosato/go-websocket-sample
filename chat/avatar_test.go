package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	gomniauthtest "github.com/stretchr/gomniauth/test"
	"github.com/stretchr/testify/assert"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	testUser := &gomniauthtest.TestUser{}
	testUser.On("AvatarURL").Return("", ErrNoAvatar)
	testChatUser := &chatUser{
		User: testUser,
	}
	url, err := authAvatar.GetAvatarURL(testChatUser)
	assert.Equal(t, err, ErrNoAvatar)

	testURL := "http://url-to-avatar/"
	testUser = &gomniauthtest.TestUser{}
	testChatUser.User = testUser
	testUser.On("AvatarURL").Return(testURL, nil)
	url, err = authAvatar.GetAvatarURL(testChatUser)
	assert.Equal(t, testURL, url)
	assert.NoError(t, err)

}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	user := &chatUser{uniqueID: "abc"}
	url, err := gravatarAvatar.GetAvatarURL(user)
	assert.NoError(t, err)

	assert.Equal(t, "//www.gravatar.com/avatar/abc", url)
}

func TestFileSystemAvatar(t *testing.T) {
	filename := filepath.Join("avatars", "xyz.png")
	ioutil.WriteFile(filename, []byte{}, 0777)
	defer func() { os.Remove(filename) }()

	var fileSystemAvatar FileSystemAvatar
	user := &chatUser{uniqueID: "xyz"}
	url, err := fileSystemAvatar.GetAvatarURL(user)
	assert.NoError(t, err)
	assert.Equal(t, "/avatars/xyz.png", url)
}
