package msq

import (
	"time"
)

type QueueConfig struct {
	Name       string
	MaxRetries int64
	MessageTTL time.Duration
}

type Queue struct {
	Connection *Connection
	Config     *QueueConfig
}

func (q *Queue) Configure(config *QueueConfig) {
	q.Config = config
}

func (q *Queue) Done(event *MessengerMessage) error {
	return q.Connection.Database().Unscoped().Delete(event).Error
}

func (q *Queue) ReQueue(event *MessengerMessage) error {
	pushback := time.Now().Add(time.Millisecond * (time.Duration(1) * 100))
	return q.Connection.Database().
		Unscoped().
		Model(event).
		Updates(map[string]interface{}{
		"available_at": pushback,
	}).Error
}

func (q *Queue) Pop() (*MessengerMessage, error) {
	message := &MessengerMessage{}

	db := q.Connection.Database()

	err := db.
		Order("created_at asc").
		Where("available_at <= ?", time.Now()).
		Where("queue_name = ?", q.Config.Name).
		First(message).Error

	if err != nil {
		return message, err
	}

	db.Delete(message)
	return message, nil
}

func (q *Queue) Failed() ([]*MessengerMessage, error) {
	messages := []*MessengerMessage{}

	db := q.Connection.Database()

	err := db.Unscoped().Order("created_at desc").
		Where("queue_name = ?", q.Config.Name).
		Find(&messages).
		Error

	if err != nil {
		return messages, err
	}

	return messages, nil
}

func (q *Queue) Push(payload Payload) (*MessengerMessage, error) {
	encodedPayload, err := payload.Marshal()

	if err != nil {
		return &MessengerMessage{}, err
	}

	event := &MessengerMessage{
		QueueName: q.Config.Name,
		Body:   string(encodedPayload),
	}

	err = q.Connection.Database().Create(event).Error

	if err != nil {
		return &MessengerMessage{}, err
	}

	return event, nil
}
