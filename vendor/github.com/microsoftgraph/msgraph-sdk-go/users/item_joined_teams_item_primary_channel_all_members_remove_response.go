package users

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemJoinedTeamsItemPrimaryChannelAllMembersRemovePostResponseable instead.
type ItemJoinedTeamsItemPrimaryChannelAllMembersRemoveResponse struct {
    ItemJoinedTeamsItemPrimaryChannelAllMembersRemovePostResponse
}
// NewItemJoinedTeamsItemPrimaryChannelAllMembersRemoveResponse instantiates a new ItemJoinedTeamsItemPrimaryChannelAllMembersRemoveResponse and sets the default values.
func NewItemJoinedTeamsItemPrimaryChannelAllMembersRemoveResponse()(*ItemJoinedTeamsItemPrimaryChannelAllMembersRemoveResponse) {
    m := &ItemJoinedTeamsItemPrimaryChannelAllMembersRemoveResponse{
        ItemJoinedTeamsItemPrimaryChannelAllMembersRemovePostResponse: *NewItemJoinedTeamsItemPrimaryChannelAllMembersRemovePostResponse(),
    }
    return m
}
// CreateItemJoinedTeamsItemPrimaryChannelAllMembersRemoveResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemJoinedTeamsItemPrimaryChannelAllMembersRemoveResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemJoinedTeamsItemPrimaryChannelAllMembersRemoveResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemJoinedTeamsItemPrimaryChannelAllMembersRemovePostResponseable instead.
type ItemJoinedTeamsItemPrimaryChannelAllMembersRemoveResponseable interface {
    ItemJoinedTeamsItemPrimaryChannelAllMembersRemovePostResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}
