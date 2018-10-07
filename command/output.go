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
	_
	StateCritical
	StateUndefined
)

// outputter is an helper function that format the output messages.
type outputter func(errCode int, format string, a ...interface{}) string

// OutputFormat returns the output helper function that is responsible
// of the command output formatting and return codes selection.
func (c *BaseCommand) Outputter() (outputter, error) {
	switch c.OutputFormat {
	case "default":
		return func(errCode int, format string, a ...interface{}) string {
			return fmt.Sprintf(format, a...)
		}, nil
	case "nagios":
		return func(errCode int, format string, a ...interface{}) string {
			return fmt.Sprintf("check_vault: "+format, a...)
		}, nil
	default:
		return nil, errors.New("Unknown outputter: " + c.OutputFormat)
	}
}
