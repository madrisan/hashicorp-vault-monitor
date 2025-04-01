## 0.9.1 -- (Apr 1, 2025)

SECURITY FIXES:

 * jwt-go allows excessive memory allocation during header parsing
   Affected verions: github.com/golang-jwt/jwt/v4 < 4.5.2
   Fix: 4.5.2
   See: https://cwe.mitre.org/data/definitions/405.html

## 0.9.0 -- (Mar 31, 2025)

SECURITY FIXES:

 * Update the go dependencies to fix several security issues.

IMPROVEMENTS:

 * Better error message for the get command.
 * Add a GitHub Workflows (as a replacement for semaphoreCI).

CHANGES:

 * Documentation updates.
 * Update tests to the new Vault API.

## 0.8.6 (Nov 27, 2021)

SECURITY FIXES:

 * Update the go dependencies to fix several security issues.

## 0.8.5 (May 5, 2020)

BUG FIXES:

 * Send all the messages to *stdout* when the Nagios outputter
   (`-output=nagios`) is selected.
   This is required because, as pointed out by *unix196*, Nagios shows an
   empty output in case of warning and error messages sent to *stderr*
   (if the stderr is not redirected to stdout).

CHANGES:

 * When the Nagios outputter is selected (`-output=nagios`), the messages
   are now printed without any color.

 * Documentation updates.

IMPROVEMENTS:

 * *hashicorp-vault-monitor* now uses Go's official dependency management
   system, Go Modules, to manage dependencies.

 * Include the error message in output when reading the environment variables.
   This will help debug the issues related to environment variables loading.
   Pull Request by *maxadamo*. Thanks!

 * Travis CI now uses go 1.14.x as build target.

 * Add CircleCI and SemaphoreCI continuous integration configurations
   with go 1.13.x and 1.14.x as build targets.

## 0.8.4 (March 10, 2020)

IMPROVEMENTS:

 * Monitor the expiration date of a Vault token via its associated
   token accessor with the new command-line option: -token-accessor.
 * Update the documentation.

## 0.8.3 (December 5, 2019)

BUG FIXES:

 * Fix (once again) the initialization of the Vault URL by ensuring that
   the command-line value has precedence over the default value and the
   VAULT_ADDR environment variable.

OTHER:

 * Travis CI: add go 1.13.x build target and remove the 1.11.x one.

## 0.8.2 (September 30, 2019)

BUG FIXES:

 * Fix the broken initialization of the Vault URL that made impossible to
   setup the Vault address via the environment variable `VAULT_ADDR`.

IMPROVEMENTS:

 * Update the documentation.
 * Add a configuration file for
   [CircleCI](https://circleci.com/gh/madrisan/hashicorp-vault-monitor).

## 0.8.1 (September 12, 2019)

IMPROVEMENTS:

 * More human readable output message for the `token-lookup` command.
   (using the time duration parser/formatter: https://github.com/hako/durafmt)

BUG FIXES:

 * The `token` switch was not available for the `token-lookup` command.
   A token could only be entered via the `VAULT_TOKEN` environment variable.

## 0.8.0 (September 11, 2019)

FEATURES:

 * New command-line check `token-lookup`

IMPROVEMENTS:

 * Update the documentation.
 * Update the test suite.
 * Rework the output module to handle warning messages.

BUG FIXES:

 * Fix all the issues reported by the golint and megacheck tools.

## 0.7.0 (August 16, 2019)

FEATURES:

 * The new command line check option `hastatus` has been added.
   This command checks the nodes status of a Vault HA Cluster.

IMPROVEMENTS:

 * Update the documentation.
 * Update the test suite.

## 0.6.2 (October 28, 2018)

IMPROVEMENTS:

 * New outputter `output=nagios`.
   This switch enables the compliance with the Nagios ouput
   (output messages and return codes).
 * Update the test suite.

## 0.6.1 (October 1, 2018)

With this first production-ready release you can monitor:

 * The status (unsealed/sealed) of a Vault server of cluster (*vault status* command)
 * Ensure a list of policies are available (*vault policies* command)
 * The read access to the Vault KV data store, both v1 and v2 (*vault get* command)
