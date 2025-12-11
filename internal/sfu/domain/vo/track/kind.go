package track

type TrackKind string

const (
	TrackKindVoice  TrackKind = "voice"
	TrackKindCamera TrackKind = "camera"
	TrackKindData   TrackKind = "data"
)
