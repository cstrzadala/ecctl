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

package testutils

import (
	"reflect"
	"testing"
)

func TestMockInitApp(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "Succeeds mocking initApp",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MockInitApp(); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("MockInitApp() error = %v, wantErr %v", err, tt.err)
				return
			}
		})
	}
}
