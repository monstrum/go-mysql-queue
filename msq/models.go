package msq

import (
	"time"
)

type MessengerMessage struct {
	ID        	uint `gorm:"primary_key"`
	CreatedAt 	time.Time
	AvailableAt time.Time
	DeliveredAt time.Time
	QueueName 	string `gorm:"type:varchar(120);index:queue_name;not null"`
	Body      	string `gorm:"type:text;body"`
	Headers     string `gorm:"type:text;headers"`
	Retries     int `gorm:"type:int;retries"`
}

// Symfony headers required
// type, the message class qualified name
// X-Message-Stamp-Symfony\\Component\\Messenger\\Stamp\\BusNameStamp, [{\"busName\":\"messenger.bus.default\"}]
// Content-Type, application\/json
type Headers map[string]interface{}

func (e *MessengerMessage) GetPayload() (Payload, error) {
	p := Payload{}
	returnPayload, err := p.UnMarshal([]byte(e.Body))

	if err != nil {
		return Payload{}, err
	}

	return *returnPayload, nil
}

func (e *MessengerMessage) GetHeaders() (Payload, error) {
	p := Payload{}
	returnPayload, err := p.UnMarshal([]byte(e.Headers))

	if err != nil {
		return Payload{}, err
	}

	return *returnPayload, nil
}
