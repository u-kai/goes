package pkg

import (
	"fmt"
	"math"
	"strconv"
	"strings"
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

func (i *ExcelIndex) Up() *ExcelIndex {
	prev := i.number.prev()
	if prev == nil {
		return nil
	}
	return &ExcelIndex{
		alphabet: i.alphabet,
		number:   prev,
	}
}
func (i *ExcelIndex) Down() *ExcelIndex {
	return &ExcelIndex{
		alphabet: i.alphabet,
		number:   i.number.next(),
	}
}
func (i *ExcelIndex) Right() *ExcelIndex {
	return &ExcelIndex{
		alphabet: i.alphabet.next(),
		number:   i.number,
	}
}
func (i *ExcelIndex) Left() *ExcelIndex {
	prev := i.alphabet.prev()
	if prev == nil {
		return nil
	}
	return &ExcelIndex{
		alphabet: prev,
		number:   i.number,
	}
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

func fromNumber(i int) *AlphabetIndex {
	relation := newAlphabetRelation()
	if relation.num >= i {
		return &AlphabetIndex{value: relation.num_to_alphabets[i]}
	}
	mod_remain := i
	mod_acc := ""

	for mod_remain > relation.num {
		mod := mod_remain % relation.num
		mod_acc = fmt.Sprintf("%s%s", relation.num_to_alphabets[mod], mod_acc)
		mod_remain = int(mod_remain / relation.num)
	}
	mod := mod_remain % relation.num
	mod_acc = fmt.Sprintf("%s%s", relation.num_to_alphabets[mod], mod_acc)
	return &AlphabetIndex{value: mod_acc}
}
func (a *AlphabetIndex) toNumber() int {
	number := 0
	split := strings.Split(a.value, "")
	for i, j := 0, len(split)-1; i < j; i, j = i+1, j-1 {
		split[i], split[j] = split[j], split[i]
	}
	relation := newAlphabetRelation()
	for digit, alphabet := range split {
		digit_effect := math.Pow(float64(relation.num), float64(digit))
		number += relation.alphabets[alphabet] * int(digit_effect)
	}
	return number
}

func (i *AlphabetIndex) next() *AlphabetIndex {
	next_number := i.toNumber() + 1
	return fromNumber(next_number)
}
func (i *AlphabetIndex) prev() *AlphabetIndex {
	next_number := i.toNumber() - 1
	if next_number < 1 {
		return nil
	}
	return fromNumber(next_number)
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

func (i *NumberIndex) prev() *NumberIndex {
	if i.value-1 < 1 {
		return nil
	}
	return &NumberIndex{value: i.value - 1}
}
func (i *NumberIndex) next() *NumberIndex {
	return &NumberIndex{value: i.value + 1}
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

type alphabetRelation struct {
	num              int
	alphabets        map[string]int
	num_to_alphabets map[int]string
}

func newAlphabetRelation() *alphabetRelation {
	ALPHABET_NUM := 26
	ALPHABETS := map[string]int{"A": 1, "B": 2, "C": 3, "D": 4, "E": 5, "F": 6, "G": 7, "H": 8, "I": 9, "J": 10, "K": 11, "L": 12, "M": 13, "N": 14, "O": 15, "P": 16, "Q": 17, "R": 18, "S": 19, "T": 20, "U": 21, "V": 22, "W": 23, "X": 24, "Y": 25, "Z": 26}
	NUM_TO_ALPHABETS := map[int]string{1: "A", 2: "B", 3: "C", 4: "D", 5: "E", 6: "F", 7: "G", 8: "H", 9: "I", 10: "J", 11: "K", 12: "L", 13: "M", 14: "N", 15: "O", 16: "P", 17: "Q", 18: "R", 19: "S", 20: "T", 21: "U", 22: "V", 23: "W", 24: "X", 25: "Y", 26: "Z"}
	return &alphabetRelation{
		num:              ALPHABET_NUM,
		alphabets:        ALPHABETS,
		num_to_alphabets: NUM_TO_ALPHABETS,
	}
}
