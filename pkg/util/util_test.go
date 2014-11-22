package util

import (
	"strings"
	"testing"
)

func TestCollapseNewlines(t *testing.T) {
	testString := "include_recipe 'nginx'\n\n\n"

	res := CollapseNewlines(testString)

	if strings.Contains(res, "\n\n\n") {
		t.Errorf("Expected resulting string not to contain 3 newlines")
	}

	if !strings.Contains(res, "\n\n") {
		t.Errorf("Expected new string to contain two newlines")
	}
}

func TestCollapseNewlines_catches_all(t *testing.T) {
	testString := "include_recipe 'nginx'\n\n\ninclude_recipe' 'apache'\n\n\n\n\npackage foo"

	res := CollapseNewlines(testString)
	if strings.Contains(res, "\n\n\n") {
		t.Errorf("Expected resulting string not to contain 3 newlines")
	}

	if strings.Contains(res, "\n\n\n\n") {
		t.Errorf("Expected resulting string to not contain 4 newlines")
	}
}

func TestIsNodeAttr(t *testing.T) {
	testMatch := "node['cookbook']['attr']"
	testNotMatch := "something[diff]"

	if !IsNodeAttr(testMatch) {
		t.Errorf("Expected %s to return true", testMatch)
	}

	if IsNodeAttr(testNotMatch) {
		t.Errorf("Expected %s to return false", testNotMatch)
	}
}

func TestFormatStrings(t *testing.T) {
	type TestStruct struct {
		NodeAttribute string
		RegString     string
	}

	tempStruct := TestStruct{
		NodeAttribute: "node['cookbook']['attr']",
		RegString:     "somestring",
	}

	FormatStrings(&tempStruct)

	if tempStruct.RegString != "'somestring'" {
		t.Errorf("Expected somestring to be changed to 'somestring'")
	}

	if tempStruct.NodeAttribute != "node['cookbook']['attr']" {
		t.Errorf("Expected a node attribute not to change")
	}
}
