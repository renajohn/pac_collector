package main

import (
	"reflect"
	"testing"
	"time"
)

type argsItem struct {
	name       string
	args       []string
	shouldFail bool
	expected   pacMonConfig
}

func TestParseCmdParams(t *testing.T) {
	expectedTable := []argsItem{
		{
			name:       "Happy case",
			args:       []string{"pacmon", "-sourceURL=ws://test", "-pollingInterval=123", "-sinkURL=http://test"},
			shouldFail: false,
			expected: pacMonConfig{
				sourceURL:       "ws://test",
				pollingInterval: time.Duration(123) * time.Second,
				sinkURL:         "http://test",
			},
		},
		{
			name:       "should error out when no source URL",
			args:       []string{"pacmon", "-pollingInterval=123", "-sinkURL=http://test"},
			shouldFail: true,
		},
		{
			name:       "should error out when no sink URL",
			args:       []string{"pacmon", "-sourceURL=http://test", "-pollingInterval=123"},
			shouldFail: true,
		},
		{
			name:       "should error out when polling interval is negative",
			args:       []string{"pacmon", "-sourceURL=ws://test", "-pollingInterval=-123", "-sinkURL=http://test"},
			shouldFail: true,
		},
		{
			name:       "should error out when polling interval is 0",
			args:       []string{"pacmon", "-sourceURL=ws://test", "-pollingInterval=0", "-sinkURL=http://test"},
			shouldFail: true,
		},
		{
			name:       "should error out when sink URL is empty",
			args:       []string{"pacmon", "-sourceURL=http://test", "-pollingInterval=123", "-sinkURL="},
			shouldFail: true,
		},
		{
			name:       "should error out when source URL is empty",
			args:       []string{"pacmon", "-sinkURL=http://test", "-pollingInterval=123", "-sourceURL="},
			shouldFail: true,
		},
		{
			name:       "should use default when no polling interval is provided",
			args:       []string{"pacmon", "-sourceURL=ws://test", "-sinkURL=http://test"},
			shouldFail: false,
			expected: pacMonConfig{
				sourceURL:       "ws://test",
				pollingInterval: time.Duration(60) * time.Second,
				sinkURL:         "http://test",
			},
		},
	}

	for _, test := range expectedTable {
		t.Run(test.name, func(t *testing.T) {
			config, err := parseCmdParams(test.args)
			if test.shouldFail && err == nil {
				t.Errorf("Following args should generate an error: %v", test.args)
			} else if !test.shouldFail && !reflect.DeepEqual(config, &test.expected) {
				t.Errorf("Expected %v, got %v", test.expected, config)
			}
		})
	}

}
