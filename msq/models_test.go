package msq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventPayload(t *testing.T) {
	message := MessengerMessage{
		Body: `{"example":"json", "payload":"4tests"}`,
	}

	assert.NotEmpty(t, message.Body)

	payload, err := message.GetPayload()

	if assert.Nil(t, err) {
		assert.Equal(t, payload["example"], "json")
		assert.Equal(t, payload["payload"], "4tests")
	}
}
