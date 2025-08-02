package main

import (
	"errors"
	"os"
	"path"
)

// ErrNoAvatarURL er feilen som returneres når
// Avatar-forekomsten ikke kan oppgi en avatar-URL.
var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL")

// Avatar representerer typer som kan representere
// brukerprofilbilder.
type Avatar interface {
	// GetAvatarURL henter avatar-URL-en for den angitte klienten,
	// eller returnerer en feilmelding hvis noe går galt.
	// ErrNoAvatarURL returneres hvis objektet ikke klarer å hente
	// en URL for den angitte klienten.
	GetAvatarURL(chatUser) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()
	if len(url) == 0 {
		return "", ErrNoAvatarURL
	}
	return url, nil
}

type GravatarAvatar struct{}

var UserGravatar GravatarAvatar

func (GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}

type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

func (FileSystemAvatar) GetAvatarURL(u chatUser) (string, error) {
	if files, err := os.ReadDir("avatars"); err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if match, _ := path.Match(u.UniqueID()+"*", file.Name()); match {
				return "/avatars/" + file.Name(), nil
			}
		}
	}
	return "", ErrNoAvatarURL
}

type TryAvatars []Avatar

func (a TryAvatars) GetAvatarURL(u ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(u); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}
