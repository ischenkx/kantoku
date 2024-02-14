package mongo

import "time"

type Resource struct {
	ID        string    `bson:"_id"`
	Data      []byte    `bson:"data"`
	Status    string    `bson:"status"`
	Version   *string   `bson:"version"`
	UpdatedAt time.Time `bson:"updated_at"`
}
