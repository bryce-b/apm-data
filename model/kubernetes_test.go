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

func TestKubernetesTransform(t *testing.T) {
	namespace, podname, poduid, nodename := "namespace", "podname", "poduid", "nodename"

	tests := []struct {
		Kubernetes Kubernetes
		Output     any
	}{
		{
			Kubernetes: Kubernetes{},
			Output:     nil,
		},
		{
			Kubernetes: Kubernetes{
				Namespace: namespace,
				NodeName:  nodename,
				PodName:   podname,
				PodUID:    poduid,
			},
			Output: map[string]any{
				"namespace": namespace,
				"node":      map[string]any{"name": nodename},
				"pod": map[string]any{
					"uid":  poduid,
					"name": podname,
				},
			},
		},
	}

	for _, test := range tests {
		output := transformAPMEvent(APMEvent{Kubernetes: test.Kubernetes})
		assert.Equal(t, test.Output, output["kubernetes"])
	}
}
