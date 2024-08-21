package structs

type Application struct {
	ID                             Snowflake                                                       `json:"id"`
	Name                           string                                                          `json:"name"`
	Icon                           *string                                                         `json:"icon,omitempty"`
	Description                    string                                                          `json:"description"`
	RPCOrigins                     []string                                                        `json:"rpc_origins,omitempty"`
	IsBotPublic                    bool                                                            `json:"bot_public"`
	IsBotRequireCodeGrant          bool                                                            `json:"bot_require_code_grant"`
	Bot                            *User                                                           `json:"bot,omitempty"`
	TermsOfServiceURL              *string                                                         `json:"terms_of_service_url,omitempty"`
	PrivacyPolicyURL               *string                                                         `json:"privacy_policy_url,omitempty"`
	Owner                          *User                                                           `json:"owner,omitempty"`
	Summary                        string                                                          `json:"summary"`
	VerifyKey                      string                                                          `json:"verify_key"`
	Team                           Team                                                            `json:"team"`
	GuildID                        *Snowflake                                                      `json:"guild_id,omitempty"`
	Guild                          *Guild                                                          `json:"guild,omitempty"`
	PrimarySKUID                   *Snowflake                                                      `json:"primary_sku_id,omitempty"`
	Slug                           *string                                                         `json:"slug,omitempty"`
	CoverImage                     *string                                                         `json:"cover_image,omitempty"`
	Flags                          Bitfield[ApplicationFlag]                                       `json:"flags"`
	ApproximateGuildCount          int                                                             `json:"approximate_guild_count"`
	RedirectURIs                   []string                                                        `json:"redirect_uris,omitempty"`
	InteractionEndpointsURL        *string                                                         `json:"interaction_endpoints_url,omitempty"`
	RoleConnectionsVerificationURL *string                                                         `json:"role_connections_verification_url,omitempty"`
	Tags                           []string                                                        `json:"tags,omitempty"`
	InstallParams                  *InstallParams                                                  `json:"install_params,omitempty"`
	IntegrationTypesConfig         map[ApplicationIntegrationType]ApplicationIntegrationTypeConfig `json:"integration_types_config,omitempty"`
	CustomInstallURL               *string                                                         `json:"custom_install_url,omitempty"`
}

type ApplicationIntegrationType int

const (
	GuildInstall ApplicationIntegrationType = 0
	UserIsntall  ApplicationIntegrationType = 1
)

type ApplicationIntegrationTypeConfig struct {
	OAuth2InstallParams *InstallParams `json:"oauth2_install_params,omitempty"`
}

type ApplicationFlag int64

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
	Scopes      []string   `json:"scopes"`
	Permissions Permission `json:"permissions"`
}

type ApplicationCommandPermissionType int

const (
	ApplicationCommandPermissionRole    ApplicationCommandPermissionType = 1
	ApplicationCommandPermissionUser    ApplicationCommandPermissionType = 2
	ApplicationCommandPermissionChannel ApplicationCommandPermissionType = 3
)

type ApplicationCommandPermissions struct {
	ID         Snowflake                        `json:"id"`
	Type       ApplicationCommandPermissionType `json:"type"`
	Permission bool                             `json:"permission"`
}
