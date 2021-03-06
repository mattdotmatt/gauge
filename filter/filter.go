package filter

import (
	"github.com/getgauge/gauge/gauge"
	"github.com/getgauge/gauge/parser"
	"github.com/getgauge/gauge/util"
)

var ExecuteTags string
var DoNotRandomize bool
var Distribute int
var NumberOfExecutionStreams int

func GetSpecsToExecute(conceptsDictionary *gauge.ConceptDictionary, specDirs []string) ([]*gauge.Specification, int) {
	specsToExecute := specsFromArgs(conceptsDictionary, specDirs)
	totalSpecs := specsToExecute
	specsToExecute = applyFilters(specsToExecute, specsFilters())
	return sortSpecsList(specsToExecute), len(totalSpecs) - len(specsToExecute)
}

func specsFilters() []specsFilter {
	return []specsFilter{&tagsFilter{ExecuteTags}, &specsGroupFilter{Distribute, NumberOfExecutionStreams}, &specRandomizer{DoNotRandomize}}
}

func applyFilters(specsToExecute []*gauge.Specification, filters []specsFilter) []*gauge.Specification {
	for _, specsFilter := range filters {
		specsToExecute = specsFilter.filter(specsToExecute)
	}
	return specsToExecute
}

func specsFromArgs(conceptDictionary *gauge.ConceptDictionary, specDirs []string) []*gauge.Specification {
	var allSpecs []*gauge.Specification
	var specs []*gauge.Specification
	var specParseResults []*parser.ParseResult
	for _, arg := range specDirs {
		specSource := arg
		if isIndexedSpec(specSource) {
			specs, specParseResults = getSpecWithScenarioIndex(specSource, conceptDictionary)
		} else {
			specs, specParseResults = parser.ParseSpecFiles(util.GetSpecFiles(specSource), conceptDictionary)
		}
		parser.HandleParseResult(specParseResults...)
		allSpecs = append(allSpecs, specs...)
	}
	return allSpecs
}

func getSpecWithScenarioIndex(specSource string, conceptDictionary *gauge.ConceptDictionary) ([]*gauge.Specification, []*parser.ParseResult) {
	specName, indexToFilter := GetIndexedSpecName(specSource)
	parsedSpecs, parseResult := parser.ParseSpecFiles(util.GetSpecFiles(specName), conceptDictionary)
	return filterSpecsItems(parsedSpecs, newScenarioIndexFilterToRetain(indexToFilter)), parseResult
}
