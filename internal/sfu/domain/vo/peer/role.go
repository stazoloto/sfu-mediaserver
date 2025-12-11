package peer

// Role - роль участника в комнате
type Role string

const (
	RoleHost     Role = "host"
	RoleListener Role = "listener"
	RoleSpeaker  Role = "speaker"
)
