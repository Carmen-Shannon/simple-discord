package structs

type OAuth2Scope string

const (
	ActivitiesReadScope                        OAuth2Scope = "activities.read"
	ActivitiesWriteScope                       OAuth2Scope = "activities.write"
	ApplicationsBuildsReadScope                OAuth2Scope = "applications.builds.read"
	ApplicationsBuildsUploadScope              OAuth2Scope = "applications.builds.upload"
	ApplicationsCommandsScope                  OAuth2Scope = "applications.commands"
	ApplicationsCommandsUpdateScope            OAuth2Scope = "applications.commands.update"
	ApplicationsCommandsPermissionsUpdateScope OAuth2Scope = "applications.commands.permissions.update"
	ApplicationsEntitlementsScope              OAuth2Scope = "applications.entitlements"
	ApplicationsStoreUpdateScope               OAuth2Scope = "applications.store.update"
	BotScope                                   OAuth2Scope = "bot"
	ConnectionsScope                           OAuth2Scope = "connections"
	DMChannelsReadScope                        OAuth2Scope = "dm_channels.read"
	EmailScope                                 OAuth2Scope = "email"
	GdmJoinScope                               OAuth2Scope = "gdm.join"
	GuildsScope                                OAuth2Scope = "guilds"
	GuildsJoinScope                            OAuth2Scope = "guilds.join"
	GuildsMembersReadScope                     OAuth2Scope = "guilds.members.read"
	IdentifyScope                              OAuth2Scope = "identify"
	MessagesReadScope                          OAuth2Scope = "messages.read"
	RelationshipsReadScope                     OAuth2Scope = "relationships.read"
	RoleConnectionsWriteScope                  OAuth2Scope = "role_connections.write"
	RPCScope                                   OAuth2Scope = "rpc"
	RPCActivitiesWriteScope                    OAuth2Scope = "rpc.activities.write"
	RPCNotificationsReadScope                  OAuth2Scope = "rpc.notifications.read"
	RPCVoiceReadScope                          OAuth2Scope = "rpc.voice.read"
	RPCVoiceWriteScope                         OAuth2Scope = "rpc.voice.write"
	VoiceScope                                 OAuth2Scope = "voice"
	WebhookIncomingScope                       OAuth2Scope = "webhook.incoming"
)
