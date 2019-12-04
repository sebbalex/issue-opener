package analyzer

import . "github.com/sebbalex/issue-opener/model"

func compareMessages(oldMessages []Message, newMessage Message) {
	// if oldMessages is empty it probably means that this
	// is the very first issue for this repo by issue-opener
	// just return newMessage
}
