package dto

import "github.com/Carmen-Shannon/simple-discord/structs"

type EditCurrentApplicationDto struct {
	CustomInstallUrl               *string                                                                          `json:"custom_install_url,omitempty"`
	Description                    *string                                                                          `json:"description,omitempty"`
	RoleConnectionsVerificationUrl *string                                                                          `json:"role_connections_verification_url,omitempty"`
	InstallParams                  *structs.InstallParams                                                           `json:"install_params,omitempty"`
	IntegrationTypesConfig         *map[structs.ApplicationIntegrationType]structs.ApplicationIntegrationTypeConfig `json:"integration_types_config,omitempty"`
	Flags                          *structs.Bitfield[structs.ApplicationFlag]                                       `json:"flags,omitempty"`
	Icon                           *string                                                                          `json:"icon,omitempty"`
	CoverImage                     *string                                                                          `json:"cover_image,omitempty"`
	InteractionsEndpointURL        *string                                                                          `json:"interaction_endpoints_url,omitempty"`
	Tags                           *[]string                                                                        `json:"tags,omitempty"`
}

type GetApplicationRoleConnectionMetadataRecordsDto struct {
	ApplicationID structs.Snowflake `json:"application_id"`
}

type GetApplicationActivityInstanceDto struct {
	ApplicationID structs.Snowflake `json:"application_id"`
	InstanceID    string            `json:"instance_id"`
}

type UpdateApplicationRoleConnectionMetadataRecordsDto struct {
	ApplicationID structs.Snowflake                           `json:"application_id"`
	Records       []structs.ApplicationRoleConnectionMetadata `json:"records"`
}
