package orm

type CollectionContract[T any] interface {
	ToJson() string
	First() T
	Empty() bool
	NotEmpty() bool
}
