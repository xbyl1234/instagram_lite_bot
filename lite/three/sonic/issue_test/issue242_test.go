/*
 * Copyright 2021 ByteDance Inc.
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
	"testing"
	"unicode/utf8"

	"github.com/bytedance/sonic"
	"github.com/stretchr/testify/require"
)

func TestIssue242_MarshalControlChars(t *testing.T) {
	for i := 0; i < utf8.RuneSelf; i++ {
		input := string([]byte{byte(i)})
		out1, err1 := sonic.ConfigStd.Marshal(input)
		out2, err2 := json.Marshal(input)
		require.NoError(t, err1)
		require.NoError(t, err2)
		require.Equal(t, out1, out2)
	}
}
