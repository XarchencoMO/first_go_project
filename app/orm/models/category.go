package models

import (
	"encoding/json"
	"first_go_project/app/orm"
	"first_go_project/app/orm/builder"
	"first_go_project/app/orm/collections"
	"log"
	"time"
)

type Category struct {
	Id        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Category) All() orm.CollectionContract[Category] {
	categories, err := builder.
		NewSqlBuilder[Category](builder.Db, "categories").
		Get()

	if err != nil {
		log.Fatal(err)
	}

	return collections.Collection[Category]{Items: categories}
}

func (Category) First() Category {
	category, err := builder.
		NewSqlBuilder[Category](builder.Db, "categories").
		First()

	if err != nil {
		log.Fatal(err)
	}

	return category
}

func (c Category) ToJson() string {
	jsonData, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Fatalf("Error converting to JSON: %v", err)
	}
	return string(jsonData)
}

func (Category) Find(id int) Category {
	category, err := builder.
		NewSqlBuilder[Category](builder.Db, "categories").
		Where("id", "=", id).
		First()

	if err != nil {
		log.Fatal(err)
	}

	return category
}
