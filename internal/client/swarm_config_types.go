// Copyright IBM Corp. 2026
// SPDX-License-Identifier: MPL-2.0

package client

// CreateSwarmConfigRequest is the request body for CreateSwarmConfig execute action.
type CreateSwarmConfigRequest struct {
	Swarm          string   `json:"swarm"`
	Name           string   `json:"name"`
	Data           string   `json:"data"`
	Labels         []string `json:"labels"`
	TemplateDriver *string  `json:"template_driver,omitempty"`
}

// RotateSwarmConfigRequest is the request body for RotateSwarmConfig execute action.
type RotateSwarmConfigRequest struct {
	Swarm  string `json:"swarm"`
	Config string `json:"config"`
	Data   string `json:"data"`
}

// RemoveSwarmConfigsRequest is the request body for RemoveSwarmConfigs execute action.
type RemoveSwarmConfigsRequest struct {
	Swarm   string   `json:"swarm"`
	Configs []string `json:"configs"`
}

// ListSwarmConfigsRequest is the request body for ListSwarmConfigs read action.
type ListSwarmConfigsRequest struct {
	Swarm string `json:"swarm"`
}

// SwarmConfigListItem represents one item from ListSwarmConfigs response.
type SwarmConfigListItem struct {
	ID   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}
