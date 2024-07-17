package common

import (
	"encoding/json"
	"time"
)

type Context map[any]any

func (self Context) Deadline() (deadline time.Time, ok bool) {
	return
}

func (self Context) Done() <-chan struct{} {
	return nil
}

func (self Context) Err() error {
	return nil
}

func (self Context) Value(key any) any {
	return self[key]
}

func (self Context) String() string {
	b, _ := json.Marshal(self)
	return string(b)
}
