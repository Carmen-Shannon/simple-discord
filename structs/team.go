package structs

type Team struct {
	Icon        *string      `json:"icon,omitempty"`
	ID          Snowflake    `json:"id"`
	Members     []TeamMember `json:"members"`
	Name        string       `json:"name"`
	OwnerUserID Snowflake    `json:"owner_user_id"`
}

type TeamMember struct {
	MembershipState MembershipState `json:"membership_state"`
	TeamID          Snowflake       `json:"team_id"`
	User            User            `json:"user"`
	Role            string          `json:"role"`
}

type MembershipState int

const (
	InvitedMembershipState  MembershipState = 1
	AcceptedMembershipState MembershipState = 2
)
