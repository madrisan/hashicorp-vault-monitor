package groups

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemTeamPrimaryChannelAllMembersAddPostResponseable instead.
type ItemTeamPrimaryChannelAllMembersAddResponse struct {
    ItemTeamPrimaryChannelAllMembersAddPostResponse
}
// NewItemTeamPrimaryChannelAllMembersAddResponse instantiates a new ItemTeamPrimaryChannelAllMembersAddResponse and sets the default values.
func NewItemTeamPrimaryChannelAllMembersAddResponse()(*ItemTeamPrimaryChannelAllMembersAddResponse) {
    m := &ItemTeamPrimaryChannelAllMembersAddResponse{
        ItemTeamPrimaryChannelAllMembersAddPostResponse: *NewItemTeamPrimaryChannelAllMembersAddPostResponse(),
    }
    return m
}
// CreateItemTeamPrimaryChannelAllMembersAddResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemTeamPrimaryChannelAllMembersAddResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemTeamPrimaryChannelAllMembersAddResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemTeamPrimaryChannelAllMembersAddPostResponseable instead.
type ItemTeamPrimaryChannelAllMembersAddResponseable interface {
    ItemTeamPrimaryChannelAllMembersAddPostResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}
