package pkg

import (
	"fmt"
	"strconv"
)

type ExcelIndex struct {
	alphabet *AlphabetIndex
	number   *NumberIndex
}

func FromStrToIndex(str string) (*ExcelIndex, error) {
	alphabet, err := newAlphabetIndex(str)
	if err != nil {
		return nil, err
	}
	number, err := newNumberIndex(str)
	if err != nil {
		return nil, err
	}
	return &ExcelIndex{alphabet: alphabet, number: number}, nil
}

type ExcelIndexError struct {
	message string
}

func (e *ExcelIndexError) Error() string {
	return e.message
}

type AlphabetIndex struct {
	value string
}

func newAlphabetIndex(index string) (*AlphabetIndex, error) {
	alphabet, err := extract_alphabet(index)
	if err != nil {
		return nil, err
	}
	return &AlphabetIndex{value: alphabet}, nil
}

func extract_alphabet(index string) (string, error) {
	if isNum(rune(index[0])) {
		return "", &ExcelIndexError{message: fmt.Sprintf("invalid index %s", index)}
	}
	result := ""
	for _, c := range index {
		if 'A' <= c && c <= 'Z' {
			result += string(c)
		}
	}
	return result, nil
}

type NumberIndex struct {
	value int
}

func (i *NumberIndex) Str() string {
	return strconv.Itoa(i.value)
}

func newNumberIndex(index string) (*NumberIndex, error) {
	num, err := extract_number(index)
	if err != nil {
		return nil, err
	}
	i, err := strconv.Atoi(num)
	if err != nil {
		return nil, err
	}
	return &NumberIndex{value: i}, nil
}
func extract_number(index string) (string, error) {
	result := ""
	for _, c := range index {
		if isNum(c) {
			result += string(c)
		}
	}
	if result == "0" {
		return "", &ExcelIndexError{message: fmt.Sprintf("invalid index %s", index)}
	}
	return result, nil
}

func isNum(c rune) bool {
	return '0' <= c && c <= '9'
}

func (i *ExcelIndex) Value() string {
	return i.alphabet.value + i.number.Str()
}
