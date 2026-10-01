package client

import "fmt"

const DefaultCredentialType = CredentialType(CredentialTypePassword)

// CredentialType is a custom type of LDAP credential.
type CredentialType int

const (
	CredentialTypeUnknown CredentialType = iota
	// CredentialTypePassword is the default for all LDAP implementations
	CredentialTypePassword
	// CredentialTypePhrase is used for to customize the RACF schema to support
	// password phrases
	CredentialTypePhrase
)

func (c CredentialType) String() string {
	switch c {
	case CredentialTypePassword:
		return "password"
	case CredentialTypePhrase:
		return "phrase"
	default:
		return "unknown"
	}
}

// SetCredentialType sets the credential type for the LDAP config given its string form.
// Returns an error if the given credential type string is unknown.
func (c *Config) SetCredentialType(credentialType string) error {
	switch credentialType {
	case CredentialTypePassword.String():
		c.CredentialType = CredentialTypePassword
	case CredentialTypePhrase.String():
		c.CredentialType = CredentialTypePhrase
	default:
		return fmt.Errorf("invalid credential_type %q", credentialType)
	}
	return nil
}
