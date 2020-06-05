package model

import "time"

type Model struct {
	// Metadata is the contents of the Chartfile.
	Metadata *Metadata `json:"metadata"`
	Path     string    `json:"path"`
	Content  []byte    `json:"content"`
	Config   []byte    `json:"content"`
}

type Metadata struct {
	Author          string            `json:"author,omitempty" yaml:"author,omitempty"`
	Created         time.Time         `json:"created,omitempty" yaml:"created,omitempty"`
	Description     string            `json:"description,omitempty" yaml:"description,omitempty"`
	Tags            []string          `json:"tags,omitempty" yaml:"tags,omitempty"`
	Labels          map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Format          string            `json:"format,omitempty" yaml:"format,omitempty"`
	Framework       string            `json:"framework,omitempty" yaml:"framework,omitempty"`
	Metrics         []Metric          `json:"metrics,omitempty" yaml:"metrics,omitempty"`
	Hyperparameters []Hyperparameter  `json:"hyperparameters,omitempty" yaml:"hyperparameters,omitempty"`
}

type Metric struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Hyperparameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
