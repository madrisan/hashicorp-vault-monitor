package groups

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemTeamPrimaryChannelAllMembersRemovePostResponseable instead.
type ItemTeamPrimaryChannelAllMembersRemoveResponse struct {
    ItemTeamPrimaryChannelAllMembersRemovePostResponse
}
// NewItemTeamPrimaryChannelAllMembersRemoveResponse instantiates a new ItemTeamPrimaryChannelAllMembersRemoveResponse and sets the default values.
func NewItemTeamPrimaryChannelAllMembersRemoveResponse()(*ItemTeamPrimaryChannelAllMembersRemoveResponse) {
    m := &ItemTeamPrimaryChannelAllMembersRemoveResponse{
        ItemTeamPrimaryChannelAllMembersRemovePostResponse: *NewItemTeamPrimaryChannelAllMembersRemovePostResponse(),
    }
    return m
}
// CreateItemTeamPrimaryChannelAllMembersRemoveResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemTeamPrimaryChannelAllMembersRemoveResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemTeamPrimaryChannelAllMembersRemoveResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemTeamPrimaryChannelAllMembersRemovePostResponseable instead.
type ItemTeamPrimaryChannelAllMembersRemoveResponseable interface {
    ItemTeamPrimaryChannelAllMembersRemovePostResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}
