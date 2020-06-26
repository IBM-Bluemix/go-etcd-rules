package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRuleManager(t *testing.T) {
	for _, erf := range []bool{true, false} {
		rm := newRuleManager(map[string]constraint{}, erf)
		rule1, err1 := NewEqualsLiteralRule("/this/is/:a/rule", nil)
		assert.NoError(t, err1)
		rm.addRule(rule1)
		rule2, err2 := NewEqualsLiteralRule("/that/is/:a/nother", nil)
		assert.NoError(t, err2)
		rm.addRule(rule2)
		rule3, err3 := NewEqualsLiteralRule("/this/is/:a", nil)
		assert.NoError(t, err3)
		rm.addRule(rule3)
		rules := rm.getStaticRules("/this/is/a/rule", nil)
		assert.Equal(t, 1, len(rules))
		for r, index := range rules {
			assert.Equal(t, 0, index)
			assert.True(t, r.satisfiable("/this/is/a/rule", nil))
		}
		rules = rm.getStaticRules("/nothing", nil)
		assert.Equal(t, 0, len(rules))
	}
}

func TestReducePrefixes(t *testing.T) {
	prefixes := map[string]int{"/servers/internal/states": 0, "/servers/internal": 1, "/servers": 1}
	prefixes = reducePrefixes(prefixes)
	assert.Equal(t, 1, len(prefixes))
	assert.Equal(t, 1, prefixes["/servers"])

}

func TestSortPrefixesByLength(t *testing.T) {
	prefixes := map[string]int{"/servers/internal": 0, "/servers/internal/states": 1, "/servers": 2}
	sorted := sortPrefixesByLength(prefixes)
	assert.Equal(t, "/servers/internal/states", sorted[2])
	assert.Equal(t, "/servers/internal", sorted[1])
	assert.Equal(t, "/servers", sorted[0])
}

func TestCombineRuleData(t *testing.T) {
	testCases := []struct {
		sourceData   [][]string
		expectedData []string
	}{
		{
			[][]string{{"/a/b/c", "/x/y/z"}, {"/a/b/c"}},
			[]string{"/a/b/c", "/x/y/z"},
		},
	}
	dummyRule := NewAndRule()
	for idx, testCase := range testCases {
		ruleIndex := 0
		source := func(_ DynamicRule) []string {
			out := testCase.sourceData[ruleIndex]
			ruleIndex++
			return out
		}
		rules := []DynamicRule{}
		for i := 0; i < len(testCase.sourceData); i++ {
			rules = append(rules, dummyRule)
		}
		compareUnorderedStringArrays(t, testCase.expectedData, combineRuleData(rules, source), "index %d", idx)
	}
}
