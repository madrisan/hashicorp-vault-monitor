## 0.8.0 (Unreleased)

FEATURES:

 * New command-line check `token-lookup`

IMPROVEMENTS:

 * Update the documentation.
 * Update the test suite.
 * Rework the output module to handle warning messages.

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
