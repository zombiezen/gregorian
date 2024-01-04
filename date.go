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
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// A Date is a Gregorian date. The zero value is January 1, year 1.
type Date struct {
	year  int
	month int
	day   int
}

// NewDate returns the Date with the given values. The arguments may be
// outside their usual ranges and will be normalized during the conversion.
func NewDate(year int, month time.Month, day int) Date {
	d := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return Date{year: d.Year() - 1, month: int(d.Month() - 1), day: d.Day() - 1}
}

// ParseDate parses a date in either ISO 8601 format (2006-01-02) or U.S. format (1/2/2006).
func ParseDate(s string) (Date, error) {
	s = strings.TrimSpace(s)
	switch {
	case s == "":
		return Date{}, errors.New("empty date")
	case strings.Contains(s, "/"):
		return parseUSDate(s)
	case strings.Contains(s, "-"):
		return parseISODate(s)
	default:
		return Date{}, fmt.Errorf("parse date %q: unknown format", s)
	}
}

func parseUSDate(s string) (Date, error) {
	switch parts := strings.Split(s, "/"); len(parts) {
	case 2:
		month, err := strconv.Atoi(parts[0])
		if err != nil {
			return Date{}, fmt.Errorf("parse US date %q: month: %v", s, err)
		}
		if !(1 <= month && month <= 12) {
			return Date{}, fmt.Errorf("parse US date %q: invalid month %d", s, month)
		}
		day, err := strconv.Atoi(parts[1])
		if err != nil {
			return Date{}, fmt.Errorf("parse US date %q: day: %v", s, err)
		}
		if !(1 <= day && day <= 31) {
			return Date{}, fmt.Errorf("parse US date %q: invalid day %d", s, day)
		}
		return NewDate(currYear(), time.Month(month), day), nil
	case 3:
		month, err := strconv.Atoi(parts[0])
		if err != nil {
			return Date{}, fmt.Errorf("parse US date %q: month: %v", s, err)
		}
		if !(1 <= month && month <= 12) {
			return Date{}, fmt.Errorf("parse US date %q: invalid month %d", s, month)
		}
		day, err := strconv.Atoi(parts[1])
		if err != nil {
			return Date{}, fmt.Errorf("parse US date %q: day: %v", s, err)
		}
		if !(1 <= day && day <= 31) {
			return Date{}, fmt.Errorf("parse US date %q: invalid day %d", s, day)
		}
		year, err := strconv.Atoi(parts[2])
		if err != nil {
			return Date{}, fmt.Errorf("parse US date %q: year: %v", s, err)
		}
		if year < 100 {
			return Date{}, fmt.Errorf("parse US date %q: short years not allowed", s)
		}
		return NewDate(year, time.Month(month), day), nil
	default:
		return Date{}, fmt.Errorf("parse US date %q: unknown format", s)
	}
}

func parseISODate(s string) (Date, error) {
	parts := strings.Split(s, "-")
	if len(parts) != 3 {
		return Date{}, fmt.Errorf("parse ISO date %q: unknown format", s)
	}
	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return Date{}, fmt.Errorf("parse ISO date %q: year: %v", s, err)
	}
	month, err := strconv.Atoi(parts[1])
	if err != nil {
		return Date{}, fmt.Errorf("parse ISO date %q: month: %v", s, err)
	}
	if !(1 <= month && month <= 12) {
		return Date{}, fmt.Errorf("parse ISO date %q: invalid month %d", s, month)
	}
	day, err := strconv.Atoi(parts[2])
	if err != nil {
		return Date{}, fmt.Errorf("parse ISO date %q: day: %v", s, err)
	}
	if !(1 <= day && day <= 31) {
		return Date{}, fmt.Errorf("parse ISO date %q: invalid day %d", s, day)
	}
	return NewDate(year, time.Month(month), day), nil
}

// Year returns the year in which d occurs.
func (d Date) Year() int {
	return d.year + 1
}

// Month returns the month of the year specified by d.
func (d Date) Month() time.Month {
	return time.Month(d.month + 1)
}

// Day returns the day of the month specified by d.
func (d Date) Day() int {
	return d.day + 1
}

// Equal reports whether d equals d2.
func (d Date) Equal(d2 Date) bool {
	return d == d2
}

// Before reports whether d is before d2.
func (d Date) Before(d2 Date) bool {
	if d.year != d2.year {
		return d.year < d2.year
	}
	if d.month != d2.month {
		return d.month < d2.month
	}
	return d.day < d2.day
}

// Add returns the date corresponding
// to adding the given number of years, months, and days to d.
func (d Date) Add(years, months, days int) Date {
	return NewDate(d.Year()+years, d.Month()+time.Month(months), d.Day()+days)
}

// IsZero reports whether d is the zero value.
func (d Date) IsZero() bool {
	return d == Date{}
}

// String returns the date in ISO 8601 format, like "2006-01-02".
func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year(), int(d.Month()), d.Day())
}

// MarshalText returns the date in ISO 8601 format, like "2006-01-02".
func (d Date) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText parses the date from ISO 8601 format, like "2006-01-02".
func (d *Date) UnmarshalText(data []byte) error {
	var err error
	*d, err = parseISODate(string(data))
	return err
}

var currYear = func() int { return time.Now().Year() }
