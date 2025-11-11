package files

import (
	"adviser-bot/lib/e"
	"adviser-bot/storage"
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

const permission = 0774

var ErrNoSaved = errors.New("No such file or directory")

type Storage struct {
	basePath string
}

func New(basePath string) Storage {
	return Storage{basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() { err = e.Wrap("can`t save page", err) }()

	fPath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(fPath, permission); err != nil {
		return err
	}

	fName, err := fileName(page)
	if err != nil {
		return err
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil
}
func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = e.Wrap("can`t pick random page", err) }()

	path := filepath.Join(s.basePath, userName)
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, ErrNoSaved
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))
	n := rand.Intn(len(files))
	file := files[n]
	return s.decodePage(filepath.Join(path, file.Name()))
}
func (s Storage) Remove(p *storage.Page) error {
	fName, err := fileName(p)
	if err != nil {
		return e.Wrap("can't remove file", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fName)

	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("can't remove file %s", path)

		return e.Wrap(msg, err)
	}

	return nil
}
func (s Storage) IsExist(p *storage.Page) (bool, error) {
	fName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("can't check if file exists", err)
	}
	path := filepath.Join(s.basePath, p.UserName, fName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file %s exists", path)
		return false, e.Wrap(msg, err)
	}

	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("can`t decode page", err)
	}
	defer file.Close()

	var p storage.Page
	if err := gob.NewEncoder(file).Encode(&p); err != nil {
		return nil, e.Wrap("can`t decode page", err)
	}
	return &p, nil
}

func fileName(page *storage.Page) (string, error) {
	return page.Hash()
}
