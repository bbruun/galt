package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"galt-etcd-client/grains"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/cameronnewman/go-flatten"
	"github.com/tidwall/gjson"
)

/*
This utility is to filter grains on the minion to make sure it only activates
on a query from the master if it full fills the search criteria.
*/
var SearchString string
var MinionGrains *grains.Grains

func Search() {

	// START SEARCH FOR GRAIN
	var searchList []string
	if strings.Contains(SearchString, " or ") {
		searchList = strings.Split(SearchString, " or ")
	} else {
		searchList = append(searchList, SearchString)
	}

	// Test output to etcd
	s := MinionGrains.ToJSON()
	fmt.Println("minonGrains: ", s)

	// Using "github.com/cameronnewman/go-flatten"
	fmt.Printf("search: %s\n", SearchString)
	if SearchString != "" {
		fmt.Println("- using \"g\" to find grain")

		flattenedObject := getFlattenMap(s)
		for searchResult := range searchList {

			search := searchList[searchResult]
			if strings.Contains(search, "=") {
				skeyval := strings.Split(search, "=")
				sreg, _ := regexp.Compile(strings.ToLower(skeyval[1]))
				for k, v := range flattenedObject {
					// fmt.Printf(" (%s = %s)\n", k, v)
					if strings.Contains(strings.ToLower(k), skeyval[0]) {
						regsearch := sreg.FindString(strings.ToLower(v))
						if regsearch != "" {
							fmt.Println(k, v)
						}
					}
				}
			} else if strings.Contains(search, ">") || strings.Contains(search, "<") {
				fmt.Println(" - numeric grain search not supported yet")
			} else {
				for k, v := range flattenedObject {
					if strings.Contains(strings.ToLower(k), strings.ToLower(search)) {
						fmt.Println(k, v)
					}
				}
			}
		}
	} else {
		fmt.Println("finding sysinfo.bios.vendor")
		vs := gjson.Get(strings.ToLower(s), "sysinfo.bios.vendor")

		fmt.Println(vs.String())
	}
	// END SEARCH FOR GRAIN
	os.Exit(0)
}
func getFlattenMap(s string) flatten.Map {
	var jsonBlob map[string]interface{}
	d := json.NewDecoder(bytes.NewReader([]byte(s)))
	d.UseNumber()
	if err := d.Decode(&jsonBlob); err != nil {
		log.Fatal(err)
	}
	// flattenedObject := flatten.Flatten(jsonBlob)
	return flatten.Flatten(jsonBlob)
}
