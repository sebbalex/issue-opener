package analyzer

import (
	"regexp"
	"strings"

	. "github.com/sebbalex/issue-opener/model"
	log "github.com/sirupsen/logrus"
)

// ParseGHComments ...
func ParseGHComments(event *Event, comments []Comment) ([]Message, error) {
	var messages []Message
	for _, c := range comments {
		m, err := parseBodyComment(c.Body)
		if err != nil {
			log.Errorf("error parsing comment: %s", err)
			return messages, nil
		}
		messages = append(messages, m)
		event.Message <- m
		log.Debugf("done appending message: %s", m)
	}
	return messages, nil
}

func parseBodyComment(body string) (Message, error) {
	// parsing body and extract validation errors
	var m Message
	var validHeader = regexp.MustCompile(`^###\ [a-zA-Z]+`)
	var validMessage = regexp.MustCompile(`^-\ [a-zA-Z]+`)
	var validFooter = regexp.MustCompile(`^##\ [a-zA-Z]+`)
	for _, line := range strings.Split(strings.TrimSuffix(body, "\r\n"), "\r\n") {
		switch {
		case validHeader.MatchString(line):
			m.Header = line
		case validMessage.MatchString(line):
			m.Message = append(m.Message, line)
		case validFooter.MatchString(line):
			m.Footer = line
		default:
			log.Debug("It doesn't match")
		}
	}
	log.Debugf("Message: %s count %d", m, len(m.Message))
	return messageToValidationErrors(&m), nil
}

func messageToValidationErrors(m *Message) Message {
	var validKey = regexp.MustCompile(`^-\ [a-zA-Z]+\ `)
	for _, mess := range m.Message {
		var e Error
		e.Key = validKey.FindString(mess)
		e.Reason = strings.Trim(strings.Replace(mess, e.Key, "", 1), " ")
		// normalizing
		e.Key = strings.Trim(strings.Replace(e.Key, "- ", "", 1), " ")
		m.ValidationErrors = append(m.ValidationErrors, e)
	}
	log.Debugf("Message: %s count %d", m.ValidationErrors, len(m.Message))
	return *m
}

func mergeMessages(messages []Message) error {
	return nil
}
