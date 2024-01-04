// Copyright 2024 Ross Light
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//		 https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package gregorian

import (
	"testing"
	"time"
)

func TestValidateDate(t *testing.T) {
	tests := []struct {
		s        string
		currYear int
		want     Date
		wantErr  bool
	}{
		{s: "2019-02-06", currYear: 2020, want: NewDate(2019, time.February, 6)},
		{s: "2019-2-6", currYear: 2020, want: NewDate(2019, time.February, 6)},
		{s: "2/6", currYear: 2020, want: NewDate(2020, time.February, 6)},
		{s: "02/06", currYear: 2020, want: NewDate(2020, time.February, 6)},
		{s: "2/6/19", currYear: 2020, wantErr: true},
		{s: "2/6/2019", currYear: 2020, want: NewDate(2019, time.February, 6)},
		{s: "2019-13-01", currYear: 2020, wantErr: true},
		{s: "2019-00-01", currYear: 2020, wantErr: true},
		{s: "2019-06-00", currYear: 2020, wantErr: true},
		{s: "2019-06-32", currYear: 2020, wantErr: true},
		{s: "13/01", currYear: 2020, wantErr: true},
		{s: "00/01", currYear: 2020, wantErr: true},
		{s: "06/00", currYear: 2020, wantErr: true},
		{s: "06/32", currYear: 2020, wantErr: true},
		{s: "13/01/2019", currYear: 2020, wantErr: true},
		{s: "00/01/2019", currYear: 2020, wantErr: true},
		{s: "06/00/2019", currYear: 2020, wantErr: true},
		{s: "06/32/2019", currYear: 2020, wantErr: true},
	}

	defer func(oldCurrYear func() int) {
		currYear = oldCurrYear
	}(currYear)
	for _, test := range tests {
		currYear = func() int { return test.currYear }
		got, err := ParseDate(test.s)
		if got != test.want || (err != nil) != test.wantErr {
			wantErr := "<nil>"
			if test.wantErr {
				wantErr = "<non-nil>"
			}
			t.Errorf("validateDate(%q) = %v, %v; want %v, %s", test.s, got, err, test.want, wantErr)
		}
	}
}
