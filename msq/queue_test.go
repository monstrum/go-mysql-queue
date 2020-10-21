package msq

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPush(t *testing.T) {
	setup()
	defer teardown()

	config := *connectionConfig
	queue, err := Connect(config)

	queue.Configure(queueConfig)

	message, err := queue.Push(payload)

	if assert.Nil(t, err) {
		assert.NotNil(t, message)

		assert.Equal(t, message.QueueName, queue.Config.Name)

		encodedPayload, err := payload.Marshal()

		if assert.Nil(t, err) {
			assert.Equal(t, message.Body, string(encodedPayload))
		}
	}
}

func TestPop(t *testing.T) {
	setup()
	defer teardown()

	config := *connectionConfig
	queue, err := Connect(config)

	queue.Configure(queueConfig)

	_, err = queue.Push(payload)

	if assert.Nil(t, err) {
		message, err := queue.Pop()

		if assert.Nil(t, err) {
			if assert.NotEqual(t, message.ID, "", "UID should not be empty") {
				encodedPayload, err := payload.Marshal()

				if assert.Nil(t, err) {
					assert.Equal(t, message.Body, string(encodedPayload), "Payload should match")
				}
			}
		}
	}
}

func TestDone(t *testing.T) {
	setup()
	defer teardown()

	config := *connectionConfig
	queue, err := Connect(config)

	queue.Configure(queueConfig)

	message, err := queue.Pop()

	if assert.Nil(t, err, "There should be an message in the queue") {
		err := queue.Done(message)

		assert.Nil(t, err, "We should be able to remove the record")

		err = connection.Database().Where("uid = ?", message.ID).First(&MessengerMessage{}).Error

		assert.NotNil(t, err, "We want the record to be missing as it should be removed")
	}
}

func TestReQueue(t *testing.T) {
	setup()
	defer teardown()

	config := *connectionConfig
	queue, err := Connect(config)

	queue.Configure(queueConfig)

	originalmessage, err := queue.Push(payload)

	if assert.Nil(t, err) {
		message, err := queue.Pop()

		if assert.Nil(t, err, "There should be an message in the queue") {
			// Simulate waiting
			time.Sleep(listenerConfig.Interval)

			err := queue.ReQueue(message)
			assert.Nil(t, err, "We should have no problem re-queuing the message")

			// Simulate waiting
			time.Sleep(listenerConfig.Interval)

			newmessage, err := queue.Pop()

			if assert.Nil(t, err, "We should find a requeued message in the queue") {
				assert.Equal(t, message.ID, newmessage.ID, "We should get back the same message")
				assert.NotEqual(t, originalmessage.Retries, newmessage.Retries, "The retries should have increased")
				assert.NotEqual(t, originalmessage.CreatedAt.UnixNano(), newmessage.CreatedAt.UnixNano(), "The created at timestamps should be different")
			}
		}
	}
}

func TestFailed(t *testing.T) {
	setup()
	defer teardown()

	config := *connectionConfig
	queue, err := Connect(config)

	queue.Configure(queueConfig)

	originalmessage, err := queue.Push(payload)

	if assert.Nil(t, err) {
		message, err := queue.Pop()

		if assert.Nil(t, err, "There should be an message in the queue") {
			// Simulate waiting
			time.Sleep(listenerConfig.Interval)

			message.Retries = 3

			err := queue.ReQueue(message)
			assert.Nil(t, err, "We should have no problem re-queuing the message")

			// Simulate waiting
			time.Sleep(listenerConfig.Interval)

			failedmessages, err := queue.Failed()

			if assert.Nil(t, err, "We should get a list of a failed messages") {
				assert.Equal(t, originalmessage.ID, failedmessages[0].ID)
				queue.Done(failedmessages[0])
			}
		}
	}
}
