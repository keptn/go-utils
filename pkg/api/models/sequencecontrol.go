package models

// SequenceControlCommand contains instructions to issue a Sequence state change request
type SequenceControlCommand struct {
	State SequenceControlState `json:"state" binding:"required"`
	Stage string               `json:"stage"`
}

type SequenceControlResponse struct {
}
