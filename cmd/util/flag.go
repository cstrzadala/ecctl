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
	"fmt"

	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/deployment/elasticsearch/instances"
	"github.com/elastic/ecctl/pkg/util"
)

// GetInstances tries to obtain a slice with the elasticsearch cluster
// instance names either from the current cluster topology or from the
// specified --instance flag in the cobra.Command.
func GetInstances(cmd *cobra.Command, params util.ClusterParams, flagName string) ([]string, error) {
	if all, _ := cmd.Flags().GetBool("all"); all {
		return instances.List(params)
	}
	return cmd.Flags().GetStringSlice(flagName)
}

// IncompatibleFlags checks if both flags have been specified, and if so
// returns an error.
func IncompatibleFlags(cmd *cobra.Command, first, second string) error {
	if cmd.Flag(first).Changed && cmd.Flag(second).Changed {
		return fmt.Errorf(
			`incompatible flags "--%s" and "--%s" specified, "--%s" will be ignored"`,
			first, second, second,
		)
	}
	return nil
}
