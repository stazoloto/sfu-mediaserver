package session

type ICECandidate struct {
	Candidate        string `json:"candidate"`
	SDPMid           string `json:"sdp_mid,omitempty"`
	SDPMLineIndex    int    `json:"sdp_mline_index,omitempty"`
	UsernameFragment string `json:"username_fragment,omitempty"`
}
