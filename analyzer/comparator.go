package analyzer

import (
	. "github.com/sebbalex/issue-opener/model"
	log "github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
)

// CompareMessages if oldMessages is empty it probably means that this
// is the very first issue for this repo by issue-opener
// just return newMessage
func CompareMessages(event *Event) error {
	var m Message
	if len(event.Message) > 0 {
		// check every message and its body and compare them with
		// event validation errors, if same, skip
		// if not delete all message attached and create
		// a new one with delta

		// code using message instead Error struct
		// var aggr []string
		// for _, m := range event.Message {
		// 	aggr = append(aggr, joinKeyValueValidationErrors(&m.ValidationErrors)...)
		// }
		// delta := deltaMessage(aggr, joinKeyValueValidationErrors(&event.ValidationError))

		var aggr []Error
		for _, m := range event.Message {
			aggr = append(aggr, m.ValidationErrors...)
		}
		delta := deltaValidationErrors(aggr, event.ValidationError)
		m.ValidationErrors = delta
		m.Append = true
		log.Debugf("delta is %v", delta)

		if len(delta) == 0 {
			event.Message = []Message{}
			return nil
		}
	} else {
		// No message exists, parse validationErrors
		// and create new issue
		// event.ValidationError to Message
		// insert back in Event
		// create issue
		m.ValidationErrors = event.ValidationError
		m.Append = false
	}
	m.Template()
	event.Message = append([]Message{}, m)
	log.Tracef("m: %s", m.String())
	log.Tracef("event: %v", event)
	return nil
}

func deltaMessage(mess []string, validationErrors []string) (delta []string) {
	log.Debugf("ve %v mess %v", validationErrors, mess)
	delta = funk.FilterString(validationErrors, func(ve string) bool {
		return !funk.ContainsString(mess, ve)
	})
	return delta
}

func deltaValidationErrors(mess []Error, validationErrors []Error) []Error {
	log.Debugf("ve %v mess %v", validationErrors, mess)
	for idx, ve := range validationErrors {
		for _, me := range mess {
			if len(validationErrors) > 0 && (ve.Key == me.Key && ve.Reason == me.Reason) {
				validationErrors = deltaValidationErrors(mess, remove(validationErrors, idx))
				return validationErrors
			}
		}
	}
	return validationErrors
}

func remove(slice []Error, s int) []Error {
	return append(slice[:s], slice[s+1:]...)
}
