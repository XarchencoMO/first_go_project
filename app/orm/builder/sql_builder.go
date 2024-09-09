package builder

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type SqlBuilder[T any] struct {
	db    *sql.DB
	table string
	where []string
	args  []interface{}
	limit int
}

// Конструктор для создания нового Builder
func NewSqlBuilder[T any](db *sql.DB, table string) *SqlBuilder[T] {
	return &SqlBuilder[T]{db: db, table: table}
}

// Реализация метода Where с использованием параметров
func (b *SqlBuilder[T]) Where(field string, operator string, value interface{}) *SqlBuilder[T] {
	condition := fmt.Sprintf(`"%s" %s $%d`, field, operator, len(b.args)+1)
	b.where = append(b.where, condition)
	b.args = append(b.args, value)
	return b
}

// Реализация метода Limit
func (b *SqlBuilder[T]) Limit(limit int) *SqlBuilder[T] {
	b.limit = limit
	return b
}

// Построение SQL-запроса
func (b *SqlBuilder[T]) buildQuery() string {
	whereClause := ""
	if len(b.where) > 0 {
		whereClause = fmt.Sprintf("WHERE %s", strings.Join(b.where, " AND "))
	}
	query := fmt.Sprintf(`SELECT * FROM "%s" %s`, b.table, whereClause)
	if b.limit > 0 {
		query = fmt.Sprintf("%s LIMIT %d", query, b.limit)
	}
	return query
}

// Универсальная функция для динамического заполнения объекта
func fillStructFromRows[T any](rows *sql.Rows, obj *T) error {
	// Получаем значение и тип структуры через рефлексию
	val := reflect.ValueOf(obj).Elem() // Получаем значение указателя на T
	modelType := val.Type()            // Получаем тип структуры

	// Создаём слайс для сканирования данных
	fields := make([]interface{}, 0)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := modelType.Field(i)

		// Проверяем, является ли поле реляционным
		if tag := fieldType.Tag.Get("relation"); tag == "true" {
			// Игнорируем реляционные поля
			continue
		}
		// Добавляем указатель на поле для сканирования
		fields = append(fields, field.Addr().Interface())
	}

	if err := rows.Scan(fields...); err != nil {
		return err
	}

	return nil
}

// Получение всех данных
func (b *SqlBuilder[T]) Get() ([]T, error) {
	// Построение SQL-запроса
	query := b.buildQuery()

	// Выполнение параметризованного запроса
	rows, err := b.db.Query(query, b.args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Приведение данных к типу T
	var results []T
	for rows.Next() {
		var item T

		err := fillStructFromRows(rows, &item)

		if err != nil {
			return nil, err
		}

		results = append(results, item)
	}

	return results, nil
}

// Получение первого элемента
func (b *SqlBuilder[T]) First() (T, error) {
	b.Limit(1)
	results, err := b.Get()
	if err != nil || len(results) == 0 {
		return *new(T), err
	}
	return results[0], nil
}
