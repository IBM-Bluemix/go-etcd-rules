package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type testExtKeyProcessor struct {
	testKeyProcessor
	workKeys map[string]string
	workTrue map[string]string
}

func (tekp *testExtKeyProcessor) isWork(key string, value *string, r readAPI) bool {
	_, ok := tekp.workTrue[key]
	tekp.workKeys[key] = ""
	return ok
}

func TestIntCrawler(t *testing.T) {
	_, c := initV3Etcd()
	kapi := c
	kapi.Put(context.Background(), "/root/child", "val1")
	kapi.Put(context.Background(), "/root1/child", "val1")

	kp := testExtKeyProcessor{
		testKeyProcessor: newTestKeyProcessor(),
		workTrue:         map[string]string{"/root/child": ""},
		workKeys:         map[string]string{},
	}

	lgr, err := zap.NewDevelopment()
	assert.NoError(t, err)
	metrics := NewMockMetricsCollector()
	metrics.SetLogger(lgr)
	expectedRuleIDs := []string{"/root/child"}
	expectedCount := []int{1}
	expectedMethods := []string{"crawler"}

	cr := intCrawler{
		kp:       &kp,
		logger:   getTestLogger(),
		prefixes: []string{"/root", "/root1"},
		kv:       c,
		metrics:  &metrics,
	}
	cr.singleRun(getTestLogger())
	if assert.Equal(t, 1, len(kp.keys)) {
		assert.Equal(t, "/root/child", kp.keys[0])
	}

	assert.Equal(t, map[string]string{
		"/root/child":  "",
		"/root1/child": "",
	}, kp.workKeys)
	assert.Equal(t, expectedRuleIDs, metrics.TimesEvaluatedRuleID)
	assert.Equal(t, expectedCount, metrics.TimesEvaluatedCount)
	assert.Equal(t, expectedMethods, metrics.TimesEvaluatedMethod)
}
