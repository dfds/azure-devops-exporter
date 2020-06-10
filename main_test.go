package main

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestConvertBuildsResponseToMapHasValidStarAndEndingsOfObjects(t *testing.T) {

	jsonFile, err := ioutil.ReadFile("test/TestConvertBuildsResponseToMapHasValidStarAndEndingsOfObjects.json")
	if err != nil {
		fmt.Print(err)
	}

	jsonString := string(jsonFile)

	stringToBuildMap := convertBuildsResponseToMap(jsonString)

	for _, value := range stringToBuildMap {
		firstChar := value[0:1]
		if firstChar != "{" {
			t.Error("First character should be: '{' but was: '" + firstChar + "'")
		}
		lastChar := value[len(value)-1:]

		if lastChar != "}" {
			t.Error("Last character should be: '}' but was: '" + lastChar + "'")
		}
	}
}
func TestRemoveExistingBuildsDoesAsStated(t *testing.T) {

	jsonFile, err := ioutil.ReadFile("test/TestRemoveExistingBuildsDoesAsStated.json")
	if err != nil {
		fmt.Print(err)
	}

	jsonString := string(jsonFile)

	idToBuildMap := convertBuildsResponseToMap(jsonString)

	existingBuildIDs := []string{"211146", "181970"}

	resultBuilds := removeExistingBuilds(existingBuildIDs, idToBuildMap)

	for _, existingBuildID := range existingBuildIDs {
		if _, ok := resultBuilds[existingBuildID]; ok {
			t.Error("Key: '" + existingBuildID + "' should not exist in resultBuilds")
		}

	}
}
