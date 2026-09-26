## 0.9.3 -- (Sep 26, 2026)

SECURITY FIXES:

 * hashicorp/vault: multiple vulnerabilities, including a critical one
   Affected versions: github.com/hashicorp/vault >= 0.3.0 < 1.19.3
   Fix: 1.19.5
   See: https://nvd.nist.gov/vuln/detail/CVE-2025-6000 (critical),
        CVE-2025-5999, CVE-2025-6203, CVE-2025-11621, CVE-2025-4656,
        and 8 further advisories reported against this range

 * grpc-go: multiple vulnerabilities, including a critical one
   Affected versions: google.golang.org/grpc < 1.79.3
   Fix: 1.79.3
   See: https://nvd.nist.gov/vuln/detail/CVE-2026-33186 (critical),
        GHSA-hrxh-6v49-42gf, CVE-2026-84304, CVE-2026-84445, CVE-2026-84303

 * dvsekhvalnov/jose2go: vulnerability in JWE/JWS handling
   Affected versions: github.com/dvsekhvalnov/jose2go < 1.7.0
   Fix: 1.7.0
   See: https://nvd.nist.gov/vuln/detail/CVE-2025-63811

 * go-jose/go-jose/v3: vulnerability in JOSE handling
   Affected versions: github.com/go-jose/go-jose/v3 < 3.0.5
   Fix: 3.0.5
   See: https://nvd.nist.gov/vuln/detail/CVE-2026-34986

 * microsoft/kiota-http-go: vulnerability in the HTTP request adapter
   Affected versions: github.com/microsoft/kiota-http-go < 1.5.5
   Fix: 1.5.5
   See: https://nvd.nist.gov/vuln/detail/CVE-2026-44503

 * mongo-driver: vulnerability in the MongoDB Go driver
   Affected versions: go.mongodb.org/mongo-driver < 1.17.7
   Fix: 1.17.7
   See: https://nvd.nist.gov/vuln/detail/CVE-2026-2303

 * Azure/go-ntlmssp: vulnerability in NTLM authentication handling
   Affected versions: github.com/Azure/go-ntlmssp < 0.1.1
   Fix: 0.1.1
   See: https://nvd.nist.gov/vuln/detail/CVE-2026-32952

 * snowflakedb/gosnowflake: low-severity vulnerability
   Affected versions: github.com/snowflakedb/gosnowflake >= 1.7.0 < 1.13.3
   Fix: 1.14.0
   See: https://nvd.nist.gov/vuln/detail/CVE-2025-46327

 * cloudflare/circl: low-severity vulnerabilities
   Affected versions: github.com/cloudflare/circl < 1.6.1
   Fix: 1.6.1
   See: https://nvd.nist.gov/vuln/detail/CVE-2025-8556,
        CVE-2026-1229

 * moby/moby (docker/docker): low-severity vulnerability
   Affected versions: github.com/docker/docker >= 26.0.0-rc1 < 28.0.0
   Fix: 28.0.0
   See: https://nvd.nist.gov/vuln/detail/CVE-2025-54410

 * filippo.io/edwards25519: low-severity vulnerability
   Affected versions: filippo.io/edwards25519 < 1.1.1
   Fix: 1.1.1
   See: https://nvd.nist.gov/vuln/detail/CVE-2026-26958

 * aws-sdk-go-v2 (service/s3 and aws/protocol/eventstream): moderate-severity
   vulnerability
   Affected versions: github.com/aws/aws-sdk-go-v2/service/s3 < 1.97.3,
     github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream < 1.7.8
   Fix: 1.97.3 and 1.7.8 respectively
   See: https://github.com/advisories/GHSA-xmrv-pmrh-hhx2

 * gorilla/websocket: moderate-severity vulnerability
   Affected versions: github.com/gorilla/websocket < 1.5.3
   Fix: 1.5.3
   See: https://github.com/advisories/GHSA-w67g-5rqw-f597

KNOWN ISSUES:

 * github.com/jackc/pgproto3/v2 and github.com/jackc/pgx/v4 are still
   flagged by GitHub's dependency scan, but neither has a fixed release in
   its major version line (v2.3.3 and v4.18.3 respectively are already the
   latest); the fix requires migrating to pgx/v5. `govulncheck` confirms
   neither is reachable from this project's own code.

IMPROVEMENTS:

 * Drop a 50+ MB kiota codegen reference file
   (`vendor/github.com/microsoftgraph/msgraph-sdk-go/kiota-dom-export.txt`)
   that isn't needed to build or test, and was tripping GitHub's large-file
   push warning.

## 0.9.2 -- (Sep 15, 2026)

SECURITY FIXES:

 * golang.org/x/text: infinite loop on invalid input
   Affected versions: golang.org/x/text < 0.39.0
   Fix: 0.41.0
   See: https://pkg.go.dev/vuln/GO-2026-5970

 * golang.org/x/net: idna package fails to reject ASCII-only Punycode-encoded
   labels
   Affected versions: golang.org/x/net < 0.55.0
   Fix: 0.57.0
   See: https://pkg.go.dev/vuln/GO-2026-5026

 * golang.org/x/net: infinite loop in the HTTP/2 transport on a bad
   SETTINGS_MAX_FRAME_SIZE
   Affected versions: golang.org/x/net < 0.53.0
   Fix: 0.57.0
   See: https://pkg.go.dev/vuln/GO-2026-4918

 * go-jose/go-jose/v4: panics during JWE decryption
   Affected versions: github.com/go-jose/go-jose/v4 < 4.1.4
   Fix: 4.1.4
   See: https://pkg.go.dev/vuln/GO-2026-4945

 * Update golang.org/x/crypto and golang.org/x/sys to their latest patched
   releases, addressing several further advisories reported against
   modules required by the project.

IMPROVEMENTS:

 * Bump the Go toolchain to 1.26 (required by the dependency updates above)
   and the `golangci-lint` version used by `make lint`/CI accordingly.

## 0.9.1 -- (Apr 2, 2025)

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
