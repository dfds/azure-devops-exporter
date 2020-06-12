package main

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestRemoveWrapperObjectHasValidStarAndEndingsOfObjects(t *testing.T) {

	jsonFile, err := ioutil.ReadFile("test/TestConvertBuildsResponseToMapHasValidStarAndEndingsOfObjects.json")
	if err != nil {
		fmt.Print(err)
	}

	jsonString := string(jsonFile)

	buildsString := removeWrapperObject(jsonString)

	firstChar := buildsString[0:1]
	if firstChar != "{" {
		t.Error("First character should be: '{' but was: '" + firstChar + "'")
	}
	lastChar := buildsString[len(buildsString)-1:]

	if lastChar != "}" {
		t.Error("Last character should be: '}' but was: '" + lastChar + "'")
	}

}
