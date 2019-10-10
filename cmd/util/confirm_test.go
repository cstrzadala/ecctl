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

package cmdutil

import (
	"reflect"
	"testing"
)

func TestActionConfirm(t *testing.T) {
	type args struct {
		actionRaw, msg string
	}
	tests := []struct {
		name    string
		args    args
		want    *bool
		wantErr bool
		err     error
	}{
		{
			name: "Does nothing if no unsafe flags have been set",
			args: args{
				actionRaw: "",
				msg:       "some message",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ActionConfirm(tt.args.actionRaw, tt.args.msg)

			if (err != nil) != tt.wantErr {
				t.Errorf("ActionConfirm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.err == nil {
				t.Errorf("ActionConfirm() expected errors = '%v' but no errors returned", tt.err)
			}

			if tt.wantErr && err.Error() != tt.err.Error() {
				t.Errorf("ActionConfirm() expected errors = '%v' but got %v", tt.err, err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ActionConfirm() expected = %v, but got %v", tt.want, got)
			}
		})
	}
}
