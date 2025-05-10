package constants

const (
	DeploymentStatusReplicaFailure = "ReplicaFailure"
	DeploymentStatusFailed         = "Failed"
	DeploymentStatusComplete       = "Complete"
	DeploymentStatusAvailable      = "Available"
	DeploymentStatusScaling        = "Scaling"
	DeploymentStatusStarting       = "Starting"
	DeploymentStatusUpdating       = "Updating"
	DeploymentStatusProgressing    = "Progressing"
	DeploymentStatusUnknown        = "Unknown"
)

// Pod Status Constants
const (
	PodStatusPending    = "Pending"
	PodStatusRunning    = "Running"
	PodStatusSucceeded  = "Succeeded"
	PodStatusFailed     = "Failed"
	PodStatusUnknown    = "Unknown"
	PodStatusStarting   = "Starting"
	PodStatusWaiting    = "Waiting"
	PodStatusTerminated = "Terminated"
)

// Pod Phases
const (
	PodPhasePending   = "Pending"
	PodPhaseRunning   = "Running"
	PodPhaseSucceeded = "Succeeded"
	PodPhaseFailed    = "Failed"
	PodPhaseUnknown   = "Unknown"
)

const (
	JobStatusInProgress = "in_progress"
	JobStatusFailed     = "failed"
	JobStatusSuccess    = "success"
	JobStatusCompleted  = "completed"
)

type AutoDeployStatus string

const (
	AutoDeployStatusInProgress AutoDeployStatus = "in_progress"
	AutoDeployStatusCompleted  AutoDeployStatus = "completed"
	AutoDeployStatusFailed     AutoDeployStatus = "failed"
	AutoDeployStatusSuccess    AutoDeployStatus = "success"
	AutoDeployStatusWaiting    AutoDeployStatus = "waiting"
)
