package structs

type Permission int64

const (
	CreateInstantInvite              Permission = 1 << 0
	KickMembers                      Permission = 1 << 1
	BanMembers                       Permission = 1 << 2
	Administrator                    Permission = 1 << 3
	ManageChannels                   Permission = 1 << 4
	ManageGuild                      Permission = 1 << 5
	AddReactions                     Permission = 1 << 6
	ViewAuditLog                     Permission = 1 << 7
	PrioritySpeaker                  Permission = 1 << 8
	Stream                           Permission = 1 << 9
	ViewChannel                      Permission = 1 << 10
	SendMessages                     Permission = 1 << 11
	SendTTSMessage                   Permission = 1 << 12
	ManageMessages                   Permission = 1 << 13
	EmbedLinks                       Permission = 1 << 14
	AttachFiled                      Permission = 1 << 15
	ReadMessageHistory               Permission = 1 << 16
	MentionEveryone                  Permission = 1 << 17
	UseExternalEmojis                Permission = 1 << 18
	ViewGuildInsights                Permission = 1 << 19
	Connect                          Permission = 1 << 20
	Speak                            Permission = 1 << 21
	MuteMembers                      Permission = 1 << 22
	DeafenMembers                    Permission = 1 << 23
	MoveMembers                      Permission = 1 << 23
	UseVAD                           Permission = 1 << 25
	ChangeNickname                   Permission = 1 << 26
	ManageNicknames                  Permission = 1 << 27
	ManageRoles                      Permission = 1 << 28
	ManageWebhooks                   Permission = 1 << 29
	ManageGuildExpressions           Permission = 1 << 30
	UseApplicationCommands           Permission = 1 << 31
	RequestToSpeak                   Permission = 1 << 32
	ManageEvents                     Permission = 1 << 33
	ManageThreads                    Permission = 1 << 34
	CreatePublicThreads              Permission = 1 << 35
	CreatePrivateThreads             Permission = 1 << 36
	UseExternalStickers              Permission = 1 << 37
	SendMessagesInThreads            Permission = 1 << 38
	UseEmbeddedActivities            Permission = 1 << 39
	ModerateMembers                  Permission = 1 << 40
	ViewCreatorMonetizationAnalytics Permission = 1 << 41
	UseSoundboard                    Permission = 1 << 42
	CreateGuildExpressions           Permission = 1 << 43
	CreateEvents                     Permission = 1 << 44
	UseExternalSounds                Permission = 1 << 45
	SendVoiceMessages                Permission = 1 << 46
	SendPolls                        Permission = 1 << 49
	UseExternalApps                  Permission = 1 << 50
)
