package model

// Status for  and run
const (
	StatusUnsent uint32 = 0
	StatusSent   uint32 = 001

	StatusInqueue uint32 = 100

	StatusInitializing uint32 = 200
	StatusInitiated    uint32 = 201

	StatusRunning uint32 = 300

	StatusSuccess uint32 = 400

	StatusPausing uint32 = 500
	StatusPaused  uint32 = 501

	StatusCanceling uint32 = 600
	StatusCanceled  uint32 = 601

	StatusFail uint32 = 700
)
