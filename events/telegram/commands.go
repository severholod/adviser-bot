package telegram

import (
	"adviser-bot/lib/utils"
	"adviser-bot/storage"
	"errors"
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new commans '%s' from user '%s'", text, username)

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHi(chatID)
	default:
		return p.tg.SendMessage(chatID, MsgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, pageUrl string, username string) (err error) {
	defer func() {
		err = utils.WrapError("can`t do command 'save page'", err)
	}()

	page := &storage.Page{
		URL:      pageUrl,
		UserName: username,
	}

	isExist, err := p.storage.IsExist(page)
	if err != nil {
		return err
	}
	if isExist {
		return p.tg.SendMessage(chatID, MsgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, MsgSaved); err != nil {
		return err
	}

	return nil
}
func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() {
		err = utils.WrapError("can`t do command 'save random'", err)
	}()

	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSaved) {
		return err
	}
	if errors.Is(err, storage.ErrNoSaved) {
		p.tg.SendMessage(chatID, MsgNoSavedPages)
	}
	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(page)
}

func (p *Processor) sendHi(chatID int) error {
	return p.tg.SendMessage(chatID, MsgHello)
}
func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, MsgHelp)
}

func isAddCmd(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
