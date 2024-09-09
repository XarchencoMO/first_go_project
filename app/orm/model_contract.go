package orm

type Model[T any] interface {
	All() CollectionContract[T]
	First() T
	Find(id int) T
	ToJson() string
}
