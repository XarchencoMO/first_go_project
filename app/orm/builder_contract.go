package orm

type Builder[T any] interface {
	Where(field string, operator string, value interface{}) Builder[T]
	Limit(limit int) Builder[T]
	Get() ([]T, error)
	First() (T, error)
}
