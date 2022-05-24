package swift

import (
	"regexp"
	"strings"
)

var tagRegex = regexp.MustCompile(`^\:?\d{1,3}[A-Z]?\:`)
var tag5Regex = regexp.MustCompile(`^[A-Z]{3}\:`)

// Parse parse swift message into data structure
func Parse(swiftText string) *Message {
	var st = strings.ReplaceAll(swiftText, "}{", "\n")
	st = strings.ReplaceAll(st, "}", "")
	st = strings.ReplaceAll(st, "{", "\n")
	st = strings.ReplaceAll(st, "\r", "")
	st = strings.TrimSpace(st)

	var sm Message
	var myBlock, myTag, myValue string
	var values = map[string][]string{}

	for _, line := range strings.Split(st, "\n") {
		line = strings.TrimSpace(line)

		if line == "" || line == "-" {
			continue
		}

		if line == "-" {
			continue
		}
		if strings.HasPrefix(line, "5:") {
			values["5"] = append(values["5"], "")
			myBlock = "5"
			continue
		}

		var tag string

		if myBlock == "5" {
			tag = tag5Regex.FindString(line)
		} else {
			tag = tagRegex.FindString(line)
		}

		if tag != "" {
			if myTag != "" {
				var _, isExist = values[myTag]

				if !isExist {
					values[myTag] = []string{}
				}

				values[myTag] = append(values[myTag], myValue)
			}

			myTag = strings.Trim(tag, ":")
			myValue = line[len(tag):]

		} else if line != "" {
			myValue = myValue + "\n" + line
		}
	}

	if myTag != "" {
		var _, isExist = values[myTag]

		if !isExist {
			values[myTag] = []string{}
		}

		values[myTag] = append(values[myTag], myValue)
	}

	sm.Body = swiftText
	sm.Tags = values
	sm.Block1 = &Block1{}
	sm.Block1.loadValues(values["1"][0])

	var val2, exists = values["2"]
	if exists && len(val2) > 0 {
		sm.Block2 = &Block2{}
		sm.Block2.loadValues(val2[0])
	}

	return &sm
}
