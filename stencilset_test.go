package fastfood

import (
	"strings"
	"testing"
)

const testStencilSet = "tests/templatepack/stencils/database/manifest.json"

func TestNewStencilSet(t *testing.T) {
	_, err := NewStencilSet(testStencilSet)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestValidStencilSet(t *testing.T) {
	s, _ := NewStencilSet(testStencilSet)

	ok, err := s.Valid()
	if !ok {
		t.Errorf("expected %s to be a valid template, returned %v", testStencilSet, err)
	}
}

func TestMergeOpts(t *testing.T) {
	s, _ := NewStencilSet(testStencilSet)

	testOpts := map[string]string{
		"database": "testdb",
		"user":     "testuser",
	}

	res := s.MergeOpts("mysql_master", testOpts)

	if res["database"] != testOpts["database"] {
		t.Errorf("expected %s as the value for database", testOpts["database"])
	}

	if res["name"] != "mysql_master" {
		t.Errorf("expected mysql_master as the name but got %s", res["name"])
	}

	if res["openfor"] != "" {
		t.Errorf("expect an empty string when no options provided but recieved %s", res["openfor"])
	}

}

func TestHelp(t *testing.T) {
	s, _ := NewStencilSet(testStencilSet)

	htext := s.Help()
	if !strings.Contains(htext, "Default Stencil: mysql_master") {
		t.Errorf("expected generated help text to display \"Default Stencil: mysql_master\"")
	}
}
