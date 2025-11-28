package track

type Settings struct {
	Bitrate    int        `json:"bitrate"`
	Framerate  int        `json:"framerate"`
	Active     bool       `json:"active"`
	Simulcast  bool       `json:"simulcast"`
	Resolution Resolution `json:"resolution"`
}
