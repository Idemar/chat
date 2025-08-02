package main

import (
	"errors"
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
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

type GravatarAvatar struct{}

var UserGravatar GravatarAvatar

func (GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			return "//www.gravatar.com/avatar/" + useridStr, nil
		}
	}
	return "", ErrNoAvatarURL
}
