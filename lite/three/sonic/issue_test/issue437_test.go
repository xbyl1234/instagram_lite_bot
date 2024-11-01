/*
 * Copyright 2023 ByteDance Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package issue_test

import (
	"encoding/json"
	"github.com/bytedance/sonic"
	"github.com/stretchr/testify/require"
	"testing"
)

// Person document for testing
type Person struct {
	Name string `bson:"name,omitempty"`
	Age  uint   `bson:"age,omitempty"`
}

func TestIssue437(t *testing.T) {
	stdv := []Person{{Name: "hello world"}}
	sonicv := []Person{{Name: "hello world"}}

	data := []byte(`[]`)
	stde := json.Unmarshal(data, &stdv)
	sonice := sonic.ConfigStd.Unmarshal(data, &sonicv)
	require.Equal(t, stde, sonice)
	require.Equal(t, stdv, sonicv)
}
