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

func Test_不正なセルのインデックスはerrorを返す(t *testing.T) {
	index, err := pkg.FromStrToIndex("invalid")
	if index != nil {
		t.Fatalf("expected nil, but got %s", index.Value())
	}
	if err == nil {
		t.Fatalf("expected error, but got nil")
	}
}
