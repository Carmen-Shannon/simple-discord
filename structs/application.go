package structs

type Application struct {
	ID                             Snowflake
	Name                           string
	Icon                           *string
	Description                    string
	RPCOrigins                     []string
	IsBotPublic                    bool
	IsBotRequireCodeGrant          bool
	Bot                            *User
	TermsOfServiceURL              *string
	PrivacyPolicyURL               *string
	Owner                          *User
	Summary                        string
	VerifyKey                      string
	Team                           Team
	GuildID                        *Snowflake
	Guild                          *Guild
	PrimarySKUID                   *Snowflake
	Slug                           *string
	CoverImage                     *string
	Flags                          ApplicationFlag
	ApproximateGuildCount          int
	RedirectURIs                   []string
	InteractionEndpointsURL        *string
	RoleConnectionsVerificationURL *string
	Tags                           []string
	InstallParams                  *InstallParams
	IntegrationTypesConfig         map[ApplicationIntegrationType]ApplicationIntegrationTypeConfig
	CustomInstallURL               *string
}

type ApplicationIntegrationType int

const (
	GuildInstall ApplicationIntegrationType = 0
	UserIsntall  ApplicationIntegrationType = 1
)

type ApplicationIntegrationTypeConfig struct {
	OAuth2InstallParams *InstallParams
}

type ApplicationFlag int

const (
	ApplicationAutoModerationOnRuleCreateBadge ApplicationFlag = 1 << 6
	GatewayPresence                            ApplicationFlag = 1 << 12
	GatewayPresenceLimited                     ApplicationFlag = 1 << 13
	GatewayGuildMembers                        ApplicationFlag = 1 << 14
	GatewayGuildMembersLimited                 ApplicationFlag = 1 << 15
	VerificationPendingGuildLimit              ApplicationFlag = 1 << 16
	EmbeddedApplicationFlag                    ApplicationFlag = 1 << 17
	GatewayMessageContent                      ApplicationFlag = 1 << 18
	GatewayMessageContentLimited               ApplicationFlag = 1 << 19
	ApplicationCommandBadge                    ApplicationFlag = 1 << 23
)

type InstallParams struct {
	Scopes      []string
	Permissions Permission
}
