package models

// SequenceTimeout is used to signal via channel that a sequence needs to be timed out
type SequenceTimeout struct {
	KeptnContext string
	LastEvent    KeptnContextExtendedCE
}

// SequenceControlState represent the wanted state of a sequence
type SequenceControlState string

const (
	// PauseSequence represent a paused sequence
	PauseSequence SequenceControlState = "pause"

	// ResumeSequence represent a sequence that was paused and should now be resumed
	ResumeSequence SequenceControlState = "resume"

	// AbortSequence represent a sequence that needs to be aborted
	AbortSequence SequenceControlState = "abort"
)

// SequenceControl represents the wanted SequenceControlState for a certain Project Stage and Context
type SequenceControl struct {
	State        SequenceControlState
	KeptnContext string
	Stage        string
	Project      string
}
