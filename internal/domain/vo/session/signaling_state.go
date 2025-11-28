package session

type SignalingState string

const (
	SignalingStateStable             SignalingState = "stable"
	SignalingStateHaveLocalOffer     SignalingState = "have-local-offer"
	SignalingStateHaveRemoteOffer    SignalingState = "have-remote-offer"
	SignalingStateHaveLocalPranswer  SignalingState = "have-local-pranswer"
	SignalingStateHaveRemotePranswer SignalingState = "have-remote-pranswer"
	SignalingStateClosed             SignalingState = "closed"
)
