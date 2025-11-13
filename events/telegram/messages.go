package telegram

const MsgHelp = `I can save and keep you pages. Also I can offer you them to read.

In order to save the page, just send me al link to it.

In order to get a random page from your list, send me command /rnd.
Caution! After that, this page will be removed from your list!`

const MsgHello = "Hi there! 👾\n\n" + MsgHelp

const (
	MsgUnknownCommand = "Unknown command 🤔"
	MsgNoSavedPages   = "You have no saved pages 🙊"
	MsgSaved          = "Saved! 👌"
	MsgAlreadyExists  = "You have already have this page in your list 🤗"
)
