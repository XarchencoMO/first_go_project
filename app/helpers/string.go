package helpers

import (
	"strconv"
)

type String struct {
	string string
}

func StringOf(string string) *String {
	return &String{string: string}
}

func (s *String) ToInt() int {
	// Преобразование строки в int
	id, err := strconv.Atoi(s.string)
	if err != nil {

	}

	return id
}
