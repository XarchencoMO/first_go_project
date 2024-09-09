package collections

import (
	"encoding/json"
	"first_go_project/app/orm"
	"log"
)

type Collection[T orm.Model[T]] struct {
	Items []T `json:"Items,omitempty"`
}

func (c Collection[T]) ToJson() string {
	jsonData, err := json.MarshalIndent(&c.Items, "", "  ")
	if err != nil {
		log.Fatalf("Error converting to JSON: %v", err)
	}
	return string(jsonData)
}

func (c Collection[T]) First() T {
	return c.Items[0]
}

func (c Collection[T]) Empty() bool {
	return len(c.Items) == 0
}

func (c Collection[T]) NotEmpty() bool {
	return !c.Empty()
}
