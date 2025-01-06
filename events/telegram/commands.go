package telegram

import (
	"context"
	"errors"
	"log"
	"strings"

	"telegram-bot/lib/e"
	"telegram-bot/storage"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

// map[chatID]Word
var usersWhoTranslate = map[int]string{}

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s", text, username)

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	if isWordToCheck(text) {
		return p.checkTranslation(chatID, text)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, words string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()

	word1, word2 := Trim(words)
	//log.Printf("trimming word1:'%s' word2:'%s'", word1, word2)

	page := &storage.Page{
		Word:        word1,
		Translation: word2,
		UserName:    username,
	}

	isExists, err := p.storage.IsExists(context.Background(), page)
	if err != nil {
		return err
	}
	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(context.Background(), page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func Trim(words string) (string, string) {
	i := strings.Index(words, "-")
	word1 := strings.TrimSpace(words[:i])
	word2 := strings.TrimSpace(words[i+1:])
	return word1, word2
}

func (p *Processor) checkTranslation(chatID int, word string) (err error) {
	w, ok := usersWhoTranslate[chatID]
	log.Printf(w, " ", ok)
	if !ok || w == "" {
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}

	if word == w {
		usersWhoTranslate[chatID] = ""
		return p.tg.SendMessage(chatID, msgCorrect)
	} else {
		return p.tg.SendMessage(chatID, msgWrong)
	}

}

func isWordToCheck(text string) bool {
	return !strings.Contains(text, "/")
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send random", err) }()

	page, err := p.storage.PickRandom(context.Background(), username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.Word); err != nil {
		return err
	}

	usersWhoTranslate[chatID] = page.Translation
	//log.Printf("in rnd got word: '%s' translation: '%s'", page.Word, page.Translation)
	//return p.storage.Remove(context.Background(), page)
	return nil
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isNewWord(text)
}

func isNewWord(text string) bool {
	i := strings.Index(text, "-")
	return i != -1
	//first := strings.TrimSpace(text[:i])
	//second := strings.TrimSpace(text[:i])
}
