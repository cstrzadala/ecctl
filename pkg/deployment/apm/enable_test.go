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

package apm

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/output"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

func TestEnable(t *testing.T) {
	type args struct {
		params EnableParams
	}
	tests := []struct {
		name string
		args args
		want *models.ApmCrudResponse
		err  error
	}{
		{
			name: "fails due to parameter validation",
			args: args{},
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
				errors.New(`id "" is invalid`),
			}},
		},
		{
			name: "Succeeds enabling an APM cluster with no tracking",
			args: args{params: EnableParams{
				ID: "d324608c97154bdba2dff97511d40368",
				API: api.NewMock(mock.Response{Response: http.Response{
					Body: mock.NewStructBody(models.ApmCrudResponse{
						ApmID:       "86d2ec6217774eedb93ba38483141997",
						SecretToken: "E1RYkvbNPvFWVvOaNGduzTUN",
					}),
					StatusCode: 201,
				}}),
			}},
			want: &models.ApmCrudResponse{
				ApmID:       "86d2ec6217774eedb93ba38483141997",
				SecretToken: "E1RYkvbNPvFWVvOaNGduzTUN",
			},
		},
		{
			name: "Returns an error when the API returns an error",
			args: args{params: EnableParams{
				ID:  "d324608c97154bdba2dff97511d40368",
				API: api.NewMock(mock.Response{Error: errors.New("an error")}),
			}},
			err: &url.Error{
				Op:  "Post",
				URL: "https://mock-host/mock-path/clusters/apm?validate_only=false",
				Err: errors.New("an error"),
			},
		},
		{
			name: "Succeeds enabling an APM cluster with tracking",
			args: args{params: EnableParams{
				TrackParams: util.TrackParams{
					Track:         true,
					Output:        output.NewDevice(new(bytes.Buffer)),
					PollFrequency: time.Millisecond,
					MaxRetries:    1,
				},
				ID: "d324608c97154bdba2dff97511d40368",
				API: api.NewMock(util.AppendTrackResponses(mock.Response{Response: http.Response{
					Body: mock.NewStructBody(models.ApmCrudResponse{
						ApmID:       "86d2ec6217774eedb93ba38483141997",
						SecretToken: "E1RYkvbNPvFWVvOaNGduzTUN",
					}),
					StatusCode: 201,
				}})...),
			}},
			want: &models.ApmCrudResponse{
				ApmID:       "86d2ec6217774eedb93ba38483141997",
				SecretToken: "E1RYkvbNPvFWVvOaNGduzTUN",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Enable(tt.args.params)
			if !reflect.DeepEqual(tt.err, err) {
				t.Errorf("Enable() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Enable() = %v, want %v", got, tt.want)
			}
		})
	}
}
