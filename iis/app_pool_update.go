package iis

import (
	"context"
	"encoding/json"
	"fmt"
)

// appPoolUpdateRequest uses omitempty on all fields so only user-specified values are sent.
// The IIS Admin API returns 500 when receiving zero-value fields it doesn't expect.
// processModelUpdate uses omitempty on int fields to avoid sending 0 (which the API rejects)
type processModelUpdate struct {
	IdleTimeout       *int64 `json:"idle_timeout,omitempty"`
	MaxProcesses      *int64 `json:"max_processes,omitempty"`
	PingingEnabled    *bool  `json:"pinging_enabled,omitempty"`
	PingInterval      *int64 `json:"ping_interval,omitempty"`
	PingResponseTime  *int64 `json:"ping_response_time,omitempty"`
	ShutdownTimeLimit *int64 `json:"shutdown_time_limit,omitempty"`
	StartupTimeLimit  *int64 `json:"startup_time_limit,omitempty"`
	IdleTimeoutAction string `json:"idle_timeout_action,omitempty"`
}

type appPoolUpdateRequest struct {
	Name                  string                `json:"name,omitempty"`
	Status                string                `json:"status,omitempty"`
	AutoStart             *bool                 `json:"auto_start,omitempty"`
	PipelineMode          string                `json:"pipeline_mode,omitempty"`
	ManagedRuntimeVersion string                `json:"managed_runtime_version,omitempty"`
	Enable32BitWin64      *bool                 `json:"enable_32bit_win64,omitempty"`
	QueueLength           *int64                `json:"queue_length,omitempty"`
	CPU                   *CPU                  `json:"cpu,omitempty"`
	ProcessModel          *processModelUpdate   `json:"process_model,omitempty"`
	Identity              *Identity             `json:"identity,omitempty"`
	Recycling             *Recycling            `json:"recycling,omitempty"`
	RapidFailProtection   *RapidFailProtection  `json:"rapid_fail_protection,omitempty"`
	ProcessOrphaning      *ProcessOrphaning     `json:"process_orphaning,omitempty"`
}

// toUpdateRequest converts an ApplicationPool to an update request, only including non-zero fields.
func toUpdateRequest(pool ApplicationPool) appPoolUpdateRequest {
	req := appPoolUpdateRequest{
		Status:       pool.Status,
		PipelineMode: pool.PipelineMode,
	}

	if pool.AutoStart {
		req.AutoStart = &pool.AutoStart
	}
	if pool.Enable32BitWin64 {
		req.Enable32BitWin64 = &pool.Enable32BitWin64
	}
	if pool.QueueLength > 0 {
		req.QueueLength = &pool.QueueLength
	}
	if pool.ManagedRuntimeVersion != "" {
		req.ManagedRuntimeVersion = pool.ManagedRuntimeVersion
	}

	// Always send nested structs if they have any non-zero fields
	emptyCPU := CPU{}
	if pool.CPU != emptyCPU {
		req.CPU = &pool.CPU
	}

	// Build process model update with only non-zero values (API rejects 0 for interval fields)
	pm := &processModelUpdate{}
	hasPM := false
	if pool.ProcessModel.IdleTimeout > 0 {
		pm.IdleTimeout = &pool.ProcessModel.IdleTimeout
		hasPM = true
	}
	if pool.ProcessModel.MaxProcesses > 0 {
		pm.MaxProcesses = &pool.ProcessModel.MaxProcesses
		hasPM = true
	}
	if pool.ProcessModel.PingingEnabled {
		pm.PingingEnabled = &pool.ProcessModel.PingingEnabled
		hasPM = true
	}
	if pool.ProcessModel.PingInterval > 0 {
		pm.PingInterval = &pool.ProcessModel.PingInterval
		hasPM = true
	}
	if pool.ProcessModel.PingResponseTime > 0 {
		pm.PingResponseTime = &pool.ProcessModel.PingResponseTime
		hasPM = true
	}
	if pool.ProcessModel.ShutdownTimeLimit > 0 {
		pm.ShutdownTimeLimit = &pool.ProcessModel.ShutdownTimeLimit
		hasPM = true
	}
	if pool.ProcessModel.StartupTimeLimit > 0 {
		pm.StartupTimeLimit = &pool.ProcessModel.StartupTimeLimit
		hasPM = true
	}
	if pool.ProcessModel.IdleTimeoutAction != "" {
		pm.IdleTimeoutAction = pool.ProcessModel.IdleTimeoutAction
		hasPM = true
	}
	if hasPM {
		req.ProcessModel = pm
	}

	emptyID := Identity{}
	if pool.Identity != emptyID {
		req.Identity = &pool.Identity
	}

	// Only send nested structs if they were explicitly configured
	if pool.Recycling.DisableOverlappedRecycle || pool.Recycling.DisableRecycleOnConfigChange ||
		pool.Recycling.LogEvents != (LogEvents{}) ||
		pool.Recycling.PeriodicRestart.TimeInterval > 0 || pool.Recycling.PeriodicRestart.PrivateMemory > 0 ||
		pool.Recycling.PeriodicRestart.RequestLimit > 0 || pool.Recycling.PeriodicRestart.VirtualMemory > 0 {
		req.Recycling = &pool.Recycling
	}

	if pool.RapidFailProtection.Enabled || pool.RapidFailProtection.LoadBalancerCapabilities != "" ||
		pool.RapidFailProtection.MaxCrashes > 0 || pool.RapidFailProtection.Interval > 0 {
		req.RapidFailProtection = &pool.RapidFailProtection
	}

	if pool.ProcessOrphaning.Enabled || pool.ProcessOrphaning.OrphanActionExe != "" {
		req.ProcessOrphaning = &pool.ProcessOrphaning
	}

	return req
}

func (client Client) UpdateAppPool(ctx context.Context, id string, pool ApplicationPool) (*ApplicationPool, error) {
	url := fmt.Sprintf("/api/webserver/application-pools/%s", id)
	updateReq := toUpdateRequest(pool)
	res, err := httpPatch(ctx, client, url, updateReq)
	if err != nil {
		return nil, err
	}
	var updated ApplicationPool
	err = json.Unmarshal(res, &updated)
	if err != nil {
		return nil, err
	}
	return &updated, nil
}
