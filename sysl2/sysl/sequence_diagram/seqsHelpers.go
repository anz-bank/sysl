package sd

import (
	"regexp"
	"strings"

	"github.com/anz-bank/sysl/src/proto"
)

func TransformBlackBoxes(blackboxes []*sysl.Attribute) [][]string {
	bbs := [][]string{}
	for _, vals := range blackboxes {
		sub_bbs := []string{}
		for _, val := range vals.GetA().Elt {
			sub_bbs = append(sub_bbs, val.GetS())
		}
		bbs = append(bbs, sub_bbs)
	}

	return bbs
}

func ParseBlackBoxesFromArgument(blackboxFlags []string) [][]string {
	bbs := [][]string{}
	for _, blackboxFlag := range blackboxFlags {
		sub_bbs := []string{}
		sub_bbs = append(sub_bbs, strings.Split(blackboxFlag, ",")...)
		bbs = append(bbs, sub_bbs)
	}

	return bbs
}

func FindMatchItems(origin string) []string {
	re := regexp.MustCompile(`(%\(\w+\))`)
	return re.FindAllString(origin, -1)
}

func RemoveWrapper(origin string) string {
	replaced := strings.Replace(origin, "%(", "", 1)
	replaced = strings.Replace(replaced, ")", "", 1)
	return replaced
}

func RemovePercentSymbol(origin string) string {
	return strings.Replace(origin, "%", "", -1)
}

func MergeAttributes(app, edpnt map[string]*sysl.Attribute) map[string]*sysl.Attribute {
	result := make(map[string]*sysl.Attribute)
	for k, v := range app {
		result[k] = v
	}
	for k, v := range edpnt {
		result[k] = v
	}

	return result
}
