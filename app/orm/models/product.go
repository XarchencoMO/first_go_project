package models

import (
	"encoding/json"
	"first_go_project/app/orm"
	"first_go_project/app/orm/builder"
	"first_go_project/app/orm/collections"
	"log"
	"time"
)

type Product struct {
	Id          int
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	categories  collections.Collection[Category] `relation:"true"`
}

func (Product) All() orm.CollectionContract[Product] {

	products, err := builder.
		NewSqlBuilder[Product](builder.Db, "products").
		Get()

	if err != nil {
		log.Fatal(err)
	}

	return collections.Collection[Product]{Items: products}
}

func (Product) First() Product {
	product, err := builder.
		NewSqlBuilder[Product](builder.Db, "products").
		First()

	if err != nil {
		log.Fatal(err)
	}

	return product
}

func (Product) Find(id int) Product {
	product, err := builder.
		NewSqlBuilder[Product](builder.Db, "products").
		Where("id", "=", id).
		First()

	if err != nil {
		log.Fatal(err)
	}

	return product
}

func (p Product) ToJson() string {
	jsonData, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		log.Fatalf("Error converting to JSON: %v", err)
	}
	return string(jsonData)
}

func (p Product) Categories() orm.CollectionContract[Category] {

	// TODO дописать many to many
	//if p.categories.Empty() {
	//	categories, err := builder.
	//		NewSqlBuilder[Category](builder.Db, "categories").
	//		Where('')
	//}

	return p.categories
}
