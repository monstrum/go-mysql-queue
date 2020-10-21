package msq

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStartStop(t *testing.T) {
	setup()
	defer teardown()

	config := *connectionConfig
	queue, err := Connect(config)

	queue.Configure(queueConfig)

	queuedmessage, err := queue.Push(payload)

	if assert.Nil(t, err) {
		listener := &Listener{
			Queue:  *queue,
			Config: listenerConfig,
		}

		assert.Equal(t, listener.Config.Interval, listenerConfig.Interval)
		assert.Equal(t, listener.Config.Timeout, listenerConfig.Timeout)

		ctx := listener.Context()

		listener.Start(func(messages []MessengerMessage) bool {
			if len(messages) > 0 {
				assert.Equal(t, queuedmessage.ID, messages[0].ID)
			}

			return true
		}, 1)

		go func() {
			assert.True(t, listener.Running, "The listener should be running")

			time.Sleep(time.Second)
			listener.Stop()
		}()

		select {
		case <-ctx.Done():
			assert.False(t, listener.Running, "The listener should no longer be running")
		}
	}
}

func TestHandleFail(t *testing.T) {
	setup()
	defer teardown()

	config := *connectionConfig
	queue, err := Connect(config)

	queue.Configure(queueConfig)

	queuedmessage, err := queue.Push(payload)

	if assert.Nil(t, err) {
		listener := &Listener{
			Queue:  *queue,
			Config: listenerConfig,
		}

		assert.Equal(t, listener.Config.Interval, listenerConfig.Interval)
		assert.Equal(t, listener.Config.Timeout, listenerConfig.Timeout)

		ctx := listener.Context()

		listener.Start(func(messages []MessengerMessage) bool {
			if len(messages) > 0 {
				assert.Equal(t, queuedmessage.ID, messages[0].ID)
			}
			return false
		}, 1)

		go func() {
			assert.True(t, listener.Running, "The listener should be started")
			time.Sleep(2 * listenerConfig.Interval)

			failedmessages, err := queue.Failed()

			if assert.Nil(t, err, "We should get a list of failed messages back") {
				assert.Equal(t, queuedmessage.ID, failedmessages[0].ID)
				queue.Done(failedmessages[0])
			}

			listener.Stop()
		}()

		select {
		case <-ctx.Done():
			assert.False(t, listener.Running, "The listener should no longer be running")
		}
	}
}

func TestHandleTimeout(t *testing.T) {
	setup()
	defer teardown()

	config := *connectionConfig
	queue, err := Connect(config)

	queue.Configure(queueConfig)

	queuedmessage, err := queue.Push(payload)

	if assert.Nil(t, err) {
		listener := &Listener{
			Queue:  *queue,
			Config: listenerConfig,
		}

		assert.Equal(t, listener.Config.Interval, listenerConfig.Interval)
		assert.Equal(t, listener.Config.Timeout, listenerConfig.Timeout)

		ctx := listener.Context()

		listener.Start(func(messages []MessengerMessage) bool {
			time.Sleep(2 * listenerConfig.Timeout)
			return false
		}, 1)

		go func() {
			assert.True(t, listener.Running, "The listener should be started")
			time.Sleep(2 * listenerConfig.Interval)

			failedmessages, err := queue.Failed()

			if assert.Nil(t, err, "We should get a list of failed messages back") {
				assert.Equal(t, queuedmessage.ID, failedmessages[0].ID)
				queue.Done(failedmessages[0])
			}

			listener.Stop()
		}()

		select {
		case <-ctx.Done():
			assert.False(t, listener.Running, "The listener should no longer be running")
		}
	}
}
