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
	Retries     int `gorm:"type:int;retries"`
}

func (e *MessengerMessage) GetPayload() (Payload, error) {
	p := Payload{}
	returnPayload, err := p.UnMarshal([]byte(e.Body))

	if err != nil {
		return Payload{}, err
	}

	return *returnPayload, nil
}
