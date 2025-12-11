package peer

// Role - роль участника в комнате
type Role string

const (
	RoleOwner  Role = "owner"
	RoleMember Role = "member"
)
