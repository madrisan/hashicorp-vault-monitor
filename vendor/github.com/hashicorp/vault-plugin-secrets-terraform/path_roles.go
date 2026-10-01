// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfc

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	userCredentialType         = "user"
	organizationCredentialType = "organization"
	teamLegacyCredentialType   = "team_legacy"
	teamCredentialType         = "team"
)

func credentialType_Values() []string {
	return []string{
		userCredentialType,
		organizationCredentialType,
		teamLegacyCredentialType,
		teamCredentialType,
	}
}

// terraformRoleEntry is a Vault role construct that maps to TFC/TFE
type terraformRoleEntry struct {
	Name           string        `json:"name"`
	Organization   string        `json:"organization,omitempty"`
	TeamID         string        `json:"team_id,omitempty"`
	UserID         string        `json:"user_id,omitempty"`
	Description    string        `json:"description,omitempty"`
	TTL            time.Duration `json:"ttl"`
	MaxTTL         time.Duration `json:"max_ttl"`
	CredentialType string        `json:"credential_type,omitempty"`
	Token          string        `json:"token,omitempty"`
	TokenID        string        `json:"token_id,omitempty"`
}

func (r *terraformRoleEntry) toResponseData() map[string]interface{} {
	respData := map[string]interface{}{
		"name":    r.Name,
		"ttl":     r.TTL.Seconds(),
		"max_ttl": r.MaxTTL.Seconds(),
	}
	if r.Description != "" {
		respData["description"] = r.Description
	}
	if r.Organization != "" {
		respData["organization"] = r.Organization
		r.CredentialType = organizationCredentialType
	}
	if r.TeamID != "" {
		respData["team_id"] = r.TeamID
		// Default to legacy team credential type
		if r.CredentialType == "" || r.CredentialType == teamLegacyCredentialType {
			r.CredentialType = teamLegacyCredentialType
			respData["credential_type"] = teamLegacyCredentialType
		} else {
			respData["credential_type"] = teamCredentialType
		}
	}
	if r.UserID != "" {
		respData["user_id"] = r.UserID
		r.CredentialType = userCredentialType
		respData["credential_type"] = userCredentialType
	}

	return respData
}

func pathRole(b *tfBackend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "role/" + framework.GenericNameRegex("name"),

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: operationPrefixTerraformCloud,
				OperationSuffix: "role",
			},

			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeLowerCaseString,
					Description: "Name of the role",
					Required:    true,
				},
				"description": {
					Type:        framework.TypeString,
					Description: "Description of the token created by the role",
				},
				"organization": {
					Type:        framework.TypeString,
					Description: "Name of the Terraform Cloud or Enterprise organization",
				},
				"team_id": {
					Type:        framework.TypeString,
					Description: "ID of the Terraform Cloud or Enterprise team under organization (e.g., settings/teams/team-xxxxxxxxxxxxx)",
				},
				"user_id": {
					Type:        framework.TypeString,
					Description: "ID of the Terraform Cloud or Enterprise user (e.g., user-xxxxxxxxxxxxxxxx)",
				},
				"ttl": {
					Type:        framework.TypeDurationSecond,
					Description: "Default lease for generated credentials. If not set or set to 0, will use system default.",
				},
				"max_ttl": {
					Type:        framework.TypeDurationSecond,
					Description: "Maximum time for role. If not set or set to 0, will use system default.",
				},
				"credential_type": {
					Type:        framework.TypeString,
					Description: "Credential type to be used for the token. Can be either 'user', 'org', 'team', or 'team_legacy'(deprecated).",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathRolesRead,
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.pathRolesWrite,
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.pathRolesDelete,
				},
			},
			HelpSynopsis:    pathRoleHelpSynopsis,
			HelpDescription: pathRoleHelpDescription,
		},
		{
			Pattern: "role/?$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: operationPrefixTerraformCloud,
				OperationVerb:   "list",
				OperationSuffix: "roles",
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: b.pathRolesList,
				},
			},

			HelpSynopsis:    pathRoleListHelpSynopsis,
			HelpDescription: pathRoleListHelpDescription,
		},
	}
}

func (b *tfBackend) pathRolesList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List(ctx, "role/")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
}

func (b *tfBackend) pathRolesRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entry, err := b.getRole(ctx, req.Storage, d.Get("name").(string))
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: entry.toResponseData(),
	}, nil
}

func (b *tfBackend) pathRolesWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	if name == "" {
		return logical.ErrorResponse("missing role name"), nil
	}

	roleEntry, err := b.getRole(ctx, req.Storage, name)
	if err != nil {
		return nil, err
	}

	if roleEntry == nil {
		roleEntry = &terraformRoleEntry{}
	}

	roleEntry.Name = name

	// Will set if users dont. Must be set for multi-team tokens
	if credentialTypeRaw, ok := d.GetOk("credential_type"); ok {
		roleEntry.CredentialType = credentialTypeRaw.(string)
		if !strutil.StrListContains(credentialType_Values(), roleEntry.CredentialType) {
			return logical.ErrorResponse("unrecognized credential type: %s", roleEntry.CredentialType), nil
		}
	}

	if organization, ok := d.GetOk("organization"); ok {
		roleEntry.Organization = organization.(string)
		if roleEntry.CredentialType == "" {
			roleEntry.CredentialType = organizationCredentialType
		}
	}

	if userID, ok := d.GetOk("user_id"); ok {
		roleEntry.UserID = userID.(string)
		if roleEntry.CredentialType == "" {
			roleEntry.CredentialType = userCredentialType
		}
	}

	if teamID, ok := d.GetOk("team_id"); ok {
		roleEntry.TeamID = teamID.(string)
		if roleEntry.CredentialType == "" || roleEntry.CredentialType == teamLegacyCredentialType {
			roleEntry.CredentialType = teamLegacyCredentialType
		}
	}

	if description, ok := d.GetOk("description"); ok {
		roleEntry.Description = description.(string)
	}

	if roleEntry.UserID != "" && (roleEntry.Organization != "" || roleEntry.TeamID != "") {
		return logical.ErrorResponse("cannot provide a user_id in combination with organization or team_id"), nil
	}

	if roleEntry.UserID == "" && roleEntry.Organization == "" && roleEntry.TeamID == "" {
		return logical.ErrorResponse("must provide an organization name, team id, or user id"), nil
	}

	if ttlRaw, ok := d.GetOk("ttl"); ok {
		roleEntry.TTL = time.Duration(ttlRaw.(int)) * time.Second
	}

	if maxTTLRaw, ok := d.GetOk("max_ttl"); ok {
		roleEntry.MaxTTL = time.Duration(maxTTLRaw.(int)) * time.Second
	}

	if roleEntry.MaxTTL != 0 && roleEntry.TTL > roleEntry.MaxTTL {
		return logical.ErrorResponse("ttl cannot be greater than max_ttl"), nil
	}

	if roleEntry.CredentialType == teamLegacyCredentialType {
		if roleEntry.Description != "" || roleEntry.TTL != 0 || roleEntry.MaxTTL != 0 {
			return logical.ErrorResponse("cannot provide description, ttl, or max_ttl with credential_type = team_legacy, try credential_type = team."), fmt.Errorf("test error")
		}
	}

	// if we're creating a role to manage a Team or Organization, we need to
	// create the token now. User tokens will be created when credentials are
	// read.
	if roleEntry.CredentialType == organizationCredentialType || roleEntry.CredentialType == teamLegacyCredentialType {
		token, err := b.createToken(ctx, req.Storage, roleEntry)
		if err != nil {
			return nil, err
		}

		roleEntry.Token = token.Token
		roleEntry.TokenID = token.ID
	}

	if err := setRole(ctx, req.Storage, name, roleEntry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *tfBackend) pathRolesDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete(ctx, "role/"+d.Get("name").(string))
	if err != nil {
		return nil, fmt.Errorf("error deleting terraform role: %w", err)
	}

	return nil, nil
}

func setRole(ctx context.Context, s logical.Storage, name string, roleEntry *terraformRoleEntry) error {
	entry, err := logical.StorageEntryJSON("role/"+name, roleEntry)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("failed to create storage entry for role")
	}

	if err := s.Put(ctx, entry); err != nil {
		return err
	}

	return nil
}

func (b *tfBackend) getRole(ctx context.Context, s logical.Storage, name string) (*terraformRoleEntry, error) {
	if name == "" {
		return nil, fmt.Errorf("missing role name")
	}

	entry, err := s.Get(ctx, "role/"+name)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	var role terraformRoleEntry

	if err := entry.DecodeJSON(&role); err != nil {
		return nil, err
	}
	return &role, nil
}

const (
	pathRoleHelpSynopsis    = `Manages the Vault role for generating Terraform Cloud / Enterprise tokens.`
	pathRoleHelpDescription = `
This path allows you to read and write roles used to generate Terraform Cloud /
Enterprise tokens. You can configure a role to manage an organization's token, a
team's token, legacy team's token, or a user's dynamic tokens, based on the 
credential_type. The credential_type is used to determine the type of token
to be generated. The credential_type can be one of the following:
- user: A user token. 
- organization: An organization token. 
- team: A team token. This is the recommend team token credential type.
- team_legacy: A legacy team token. This is the default credential type if
  team_id is set but credential_type is left empty.

credential_type "user" can have multiple API tokens. To manage a user token, you 
can user_id and credential_type "user". When issuing a call to create creds, this role
will be used to generate the token. 

credential_type "team" can have multiple API tokens. This is the recommended 
team token credential type. To manage a team token, you can set a team_id 
and set credential_type to "team". When issuing a call to create creds, this role 
will be used to generate the token. You can set a ttl and max_ttl. Max_ttl will 
also set an expiration timer on the terraform token (including the system max ttl).

credential_type "organization" or "team_legacy" can only have one active token at a
time. When a new token is created, the old token will be revoked. This is
because Terraform Cloud/Enterprise does not support multiple active tokens for these
types.

`

	pathRoleListHelpSynopsis    = `List the existing roles in Terraform Cloud / Enterprise backend`
	pathRoleListHelpDescription = `Roles will be listed by the role name.`
)
