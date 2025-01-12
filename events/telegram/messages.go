package telegram

const msgHelp = `I can save and keep new words to learn. Also I can send a random word to translate.

In order to save the word and its translation, just send me words in 'word - translation' form.

In order to get random words from your dictionary, send me command /rnd. 

In order to stop getting words, send me command /stop.`

const msgHello = "Hi there! \n\n" + msgHelp

const (
	msgUnknownCommand = "Unknown command"
	msgNoSavedPages   = "You have no saved words"
	msgSaved          = "Saved!"
	msgAlreadyExists  = "You already have this word in your list"
	msgCorrect        = "Correct!"
	msgWrong          = "Wrong! Should be "
	msgStop           = "Stopped learning. Come back any time"
)
