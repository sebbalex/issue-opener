package analyzer

import (
	"encoding/json"
	"net/url"
	"testing"

	. "github.com/sebbalex/issue-opener/model"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var urlString string = "https://raw.githubusercontent.com/AgID/pat/master/publiccode.yml"
var valid bool = false
var valErrors string = `[
		{"key": "name", "reason": "missing mandatory key"}, 
		{"key": "localisation_ready", "reason": "missing mandatory key"}
	]`
var valErrors2 string = `[
		{"key": "name", "reason": "missing mandatory key"}, 
		{"key": "name", "reason": "too short, need 35 chars"}, 
		{"key": "generic_name", "reason": "missing mandatory key"}, 
		{"key": "publiccodeYMLVersion", "reason": "missing mandatory key"}, 
		{"key": "localisation_ready", "reason": "missing mandatory key"}
	]`
var valErrors2Delta string = `[
		{"key": "name", "reason": "too short, need 35 chars"}, 
		{"key": "generic_name", "reason": "missing mandatory key"}, 
		{"key": "publiccodeYMLVersion", "reason": "missing mandatory key"}
	]`
var valErrorsStringArray = []string{
	"name missing mandatory key",
	"localisation_ready missing mandatory key",
}

func createEvent() *Event {
	event := Event{}
	event.URL, _ = url.Parse(urlString)
	event.Valid = false
	json.Unmarshal([]byte(valErrors), &event.ValidationError)
	// event.Message = make(chan Message, 100)

	return &event
}

func fillMessages(event *Event) {
	message := Message{}
	message.ValidationErrors = event.ValidationError
	message.URL = event.URL
	message.Template()
	event.Message = append(event.Message, message)
}

func TestCompareMessages(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	event := createEvent()
	CompareMessages(event)
	m := event.Message[0]
	assert.Len(t, m.ValidationErrors, 2)
	assert.Equal(t, event.ValidationError, m.ValidationErrors)
	assert.Len(t, m.Message, 2)
}
func TestCompareMessagesWithEqualMessage(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	event := createEvent()
	fillMessages(event)
	CompareMessages(event)
	m := event.Message[0]
	assert.Len(t, m.ValidationErrors, 2)
	assert.Equal(t, event.ValidationError, m.ValidationErrors)
	assert.Len(t, m.Message, 2)
}

func TestDeltaMessage(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	delta := deltaMessage(valErrorsStringArray, valErrorsStringArray)
	assert.Len(t, delta, 0)
	assert.Empty(t, delta)

	delta = deltaMessage([]string{"hey", "how", "are", "you"}, []string{"hey", "i'm", "fine", "and", "you"})
	assert.Len(t, delta, 3)
	assert.Equal(t, []string{"i'm", "fine", "and"}, delta)
	log.Infof("delta: %v", delta)
}

func TestDeltaVE(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	event := createEvent()
	var verr []Error
	json.Unmarshal([]byte(valErrors), &verr)
	delta := deltaValidationErrors(event.ValidationError, verr)
	log.Debugf("delta is %v", delta)
	assert.Len(t, delta, 0)
	assert.Empty(t, delta)

	var verrDelta []Error
	json.Unmarshal([]byte(valErrors2), &verr)
	json.Unmarshal([]byte(valErrors2Delta), &verrDelta)
	event = createEvent()
	delta = deltaValidationErrors(event.ValidationError, verr)
	assert.Len(t, delta, 3)
	assert.Equal(t, verrDelta, delta)
	log.Infof("delta: %v and verrDelta: %v", delta, verrDelta)
}
func TestRemoveSliceElement(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	var verr []Error
	var verr2 []Error
	json.Unmarshal([]byte(valErrors), &verr)
	json.Unmarshal([]byte(valErrors2), &verr2)
	verr2 = remove(verr2, 3)
	verr2 = remove(verr2, 2)
	verr2 = remove(verr2, 1)
	assert.Len(t, verr2, 2)
	assert.Equal(t, verr, verr2)
	log.Infof("delta: %v", verr2)
}
