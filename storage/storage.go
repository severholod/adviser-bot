package storage

import (
	"adviser-bot/lib/utils"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
)

var ErrNoSaved = errors.New("No such file or directory")

type Storage interface {
	Save(p *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(p *Page) error
	IsExist(p *Page) (bool, error)
}

type Page struct {
	URL      string
	UserName string
}

func (p Page) Hash() (string, error) {
	h := sha1.New()
	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", utils.WrapError("can`t calculate hash", err)
	}
	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", utils.WrapError("can`t calculate hash", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
