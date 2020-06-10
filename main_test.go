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

	existingBuildIds := []string{"211146", "181970"}

	resultBuilds := removeExistingBuilds(existingBuildIds, idToBuildMap)

	for _, existingBuildId := range existingBuildIds {
		if _, ok := resultBuilds[existingBuildId]; ok {
			t.Error("Key: '" + existingBuildId + "' should not exist in resultBuilds")
		}

	}
}
