package structs

type Team struct {
	Icon        *string
	ID          Snowflake
	Members     []TeamMember
	Name        string
	OwnerUserID Snowflake
}

type TeamMember struct {
	MembershipState MembershipState
	TeamID          Snowflake
	User            User
	Role            string
}

type MembershipState int

const (
	InvitedMembershipState  MembershipState = 1
	AcceptedMembershipState MembershipState = 2
)
