package pkg_test

import (
	"testing"

	"github.com/u-kai/goes/pkg"
)

func Test_セルのインデックスを返す(t *testing.T) {
	index, _ := pkg.FromStrToIndex("A1")
	if index.Value() != "A1" {
		t.Fatalf("expected A1, but got %s", index.Value())
	}
}

func Test_セルからそのセルの右のセルを作成することができる(t *testing.T) {
	index, _ := pkg.FromStrToIndex("A1")
	rightIndex := index.Right()
	if rightIndex.Value() != "B1" {
		t.Fatalf("expected B1, but got %s", rightIndex.Value())
	}
}

func Test_セルからそのセルの下のセルを作成することができる(t *testing.T) {
	index, _ := pkg.FromStrToIndex("A1")
	downIndex := index.Down()
	if downIndex.Value() != "A2" {
		t.Fatalf("expected A2, but got %s", downIndex.Value())
	}
}

func Test_セルからそのセルの上のセルを作成することができる(t *testing.T) {
	index, _ := pkg.FromStrToIndex("A2")
	upIndex := index.Up()
	if upIndex.Value() != "A1" {
		t.Fatalf("expected A1, but got %s", upIndex.Value())
	}
}

func Test_セルからそのセルの左のセルを作成することができる(t *testing.T) {
	index, _ := pkg.FromStrToIndex("B1")
	leftIndex := index.Left()
	if leftIndex.Value() != "A1" {
		t.Fatalf("expected A1, but got %s", leftIndex.Value())
	}
}
func Test_1セルよりも上のセルを指定した場合はnilを返す(t *testing.T) {
	index, _ := pkg.FromStrToIndex("A1")
	upIndex := index.Up()
	if upIndex != nil {
		t.Fatalf("expected A1, but got %v", upIndex)
	}
}
func Test_Aセルよりも左のセルを指定した場合はnilを返す(t *testing.T) {
	index, _ := pkg.FromStrToIndex("A1")
	leftIndex := index.Left()
	if leftIndex != nil {
		t.Fatalf("expected A1, but got %s", leftIndex.Value())
	}
}

func Test_不正なセルのインデックスはerrorを返す(t *testing.T) {
	index, err := pkg.FromStrToIndex("invalid")
	if index != nil {
		t.Fatalf("expected nil, but got %s", index.Value())
	}
	if err == nil {
		t.Fatalf("expected error, but got nil")
	}
}
