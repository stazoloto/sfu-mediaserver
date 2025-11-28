package room

// RoomType - тип комнаты (конференция, вебинар)
type RoomType string

const (
	RoomTypeWebinar    RoomType = "webinar"    // выступление на большую аудиторию
	RoomTypeConference RoomType = "conference" // для совместной работы
)
