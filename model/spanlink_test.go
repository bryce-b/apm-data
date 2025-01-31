// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpanLinkFields(t *testing.T) {
	tests := []struct {
		Input    SpanLink
		Expected map[string]any
	}{{
		Input: SpanLink{
			Span: Span{ID: "span_id"},
		},
		Expected: map[string]any{
			"span": map[string]any{"id": "span_id"},
		},
	}, {
		Input: SpanLink{
			Trace: Trace{ID: "trace_id"},
		},
		Expected: map[string]any{
			"trace": map[string]any{"id": "trace_id"},
		},
	}, {
		Input: SpanLink{
			Span:  Span{ID: "span_id"},
			Trace: Trace{ID: "trace_id"},
		},
		Expected: map[string]any{
			"span":  map[string]any{"id": "span_id"},
			"trace": map[string]any{"id": "trace_id"},
		},
	}}
	for _, test := range tests {
		output := transformAPMEvent(APMEvent{Span: &Span{Links: []SpanLink{test.Input}}})
		assert.Equal(t, []any{test.Expected}, output["span"].(map[string]any)["links"])
	}
}
