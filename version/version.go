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

package version

import (
	"bytes"
	"fmt"
)

var (
	// GitCommit contains the git commit that was compiled.
	// This will be filled in by the compiler.
	GitCommit string

	// Version contains the version number that is being run at the moment.
	Version = "unknown"

	// VersionPrerelease contains a pre-release marker for the version.
	// If this is "" (empty string) then it means that it is a final release.
	// Otherwise, this is a pre-releas such as "dev" (in development), "beta", "rc1", etc.
	VersionPrerelease = "unknown"
)

// VersionInfo contains the HashiCorp Vault Monitor version informations.
type VersionInfo struct {
	Revision          string
	Version           string
	VersionPrerelease string
}

// GetVersion returns the VersionInfo structure.
func GetVersion() *VersionInfo {
	ver := Version
	rel := VersionPrerelease

	return &VersionInfo{
		Revision:          GitCommit,
		Version:           ver,
		VersionPrerelease: rel,
	}
}

// VersionNumber returns the current release version of HashiCorp Vault Monitor.
func (c *VersionInfo) VersionNumber() string {
	if c.Version == "unknown" && c.VersionPrerelease == "unknown" {
		return "(version unknown)"
	}

	version := c.Version

	if c.VersionPrerelease != "" {
		version = fmt.Sprintf("%s-%s", version, c.VersionPrerelease)
	}

	return version
}

// FullVersionNumber returns the full version informations of HashiCorp Vault Monitor.
func (c *VersionInfo) FullVersionNumber(rev bool) string {
	var versionString bytes.Buffer

	if c.Version == "unknown" && c.VersionPrerelease == "unknown" {
		return "HashiCorp Vault Monitor (version unknown)"
	}

	fmt.Fprintf(&versionString, "HashiCorp Vault Monitor v%s", c.Version)
	if c.VersionPrerelease != "" {
		fmt.Fprintf(&versionString, "-%s", c.VersionPrerelease)
	}

	if rev && c.Revision != "" {
		fmt.Fprintf(&versionString, " (%s)", c.Revision)
	}

	return versionString.String()
}
