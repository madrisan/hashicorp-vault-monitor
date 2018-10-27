/*
  Copyright 2018 Davide Madrisan <davide.madrisan@gmail.com>

  Licensed under the Mozilla Public License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      https://www.mozilla.org/en-US/MPL/2.0/

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package command

import (
	"errors"
	"fmt"
)

// (Nagios compatible) return codes constants.
const (
	StateOk int = iota
	StateWarning
	StateCritical
	StateUndefined
)

// Outputter holds the output functions that are monitoring tool dependent.
type Outputter struct {
	Output func(format string, a ...interface{})
	Error  func(format string, a ...interface{})
}

// OutputHandle returns the output helper function that is responsible
// of the command output formatting and return codes selection.
func (c *BaseCommand) OutputHandle() (*Outputter, error) {
	switch c.OutputFormat {
	case "default":
		return &Outputter{
			Output: func(format string, a ...interface{}) {
				c.Ui.Output(fmt.Sprintf(format, a...))
			},
			Error: func(format string, a ...interface{}) {
				c.Ui.Error(fmt.Sprintf(format, a...))
			},
		}, nil
	case "nagios":
		return &Outputter{
			Output: func(format string, a ...interface{}) {
				fmt.Printf("vault OK - "+format+"\n", a...)
			},
			Error: func(format string, a ...interface{}) {
				fmt.Printf("vault CRITICAL - "+format+"\n", a...)
			},
		}, nil
	default:
		return nil, errors.New("Unknown outputter: " + c.OutputFormat)
	}
}
