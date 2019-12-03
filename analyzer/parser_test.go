package analyzer

import (
	"github.com/sebbalex/issue-opener/model"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseGHComments(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	var header = "### Developers Italia - Issue Opener for publiccode.yml"
	var comment = "### Developers Italia - Issue Opener for publiccode.yml\r\nWe have discovered potential issue on validating your `publiccode.yml`.\r\nHere some details:\r\n\r\n- url missing mandatory key\r\n- name missing mandatory key\r\n- longDescription too short (2), min 500 chars\r\n\r\nPlease review your pubbliccode.\r\n## find out more: https://developers.italia.it"
	var footer = "## find out more: https://developers.italia.it"
	message, err := parseBodyComment(comment)
	if err != nil {
		t.Fatal("error parsing body comment")
	}
	assert.Equal(t, message.Header, header)
	assert.Equal(t, len(message.Message), 3)
	assert.Equal(t, len(message.ValidationErrors), 3)
	assert.Equal(t, message.Footer, footer)
}
func TestParseBodyComment(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	keys := []string{"url", "name", "longDescription"}
	reasons := []string{"missing mandatory key", "missing mandatory key", "too short (2), min 500 chars"}
	var errors []model.Error
	for index := 0; index < len(keys); index++ {
		errors = append(errors,
			model.Error{
				Key:    keys[index],
				Reason: reasons[index],
			},
		)

	}

	message := messageToValidationErrors(
		&Message{
			Message: []string{
				"- url missing mandatory key",
				"- name missing mandatory key",
				"- longDescription too short (2), min 500 chars",
			},
		})

	assert.Equal(t, len(message.ValidationErrors), 3)
	assert.Equal(t, message.ValidationErrors, errors)

}
