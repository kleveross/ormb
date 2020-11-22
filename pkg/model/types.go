package model

import "time"

type Model struct {
	// Metadata is the contents of the Chartfile.
	Metadata *Metadata `json:"metadata,omitempty"`
	Path     string    `json:"path,omitempty"`
	Content  []byte    `json:"content,omitempty"`
	Config   []byte    `json:"config,omitempty"`
}

type Metadata struct {
	Author      string            `json:"author,omitempty" yaml:"author,omitempty"`
	Created     time.Time         `json:"created,omitempty" yaml:"created,omitempty"`
	Description string            `json:"description,omitempty" yaml:"description,omitempty"`
	Tags        []string          `json:"tags,omitempty" yaml:"tags,omitempty"`
	Labels      map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Format      string            `json:"format,omitempty" yaml:"format,omitempty"`
	// GPUType is for TensorRT format only, it must be set when extract signature or serve
	// as a online service, otherwise， it can not extract or serve as a service.
	// for other model format, you can set empty string or not set.
	GPUType            string           `json:"gpuType,omitempty" yaml:"gpuType,omitempty"`
	Framework          string           `json:"framework,omitempty" yaml:"framework,omitempty"`
	Metrics            []Metric         `json:"metrics,omitempty" yaml:"metrics,omitempty"`
	Hyperparameters    []Hyperparameter `json:"hyperparameters,omitempty" yaml:"hyperparameters,omitempty"`
	Signature          *Signature       `json:"signature,omitempty" yaml:"signature,omitempty"`
	Training           *Training        `json:"training,omitempty" yaml:"training,omitempty"`
	Dataset            *Dataset         `json:"dataset,omitempty" yaml:"dataset,omitempty"`
	DirectoryStructure []string         `json:"directoryStructure,omitempty" yaml:"directoryStructure,omitempty"`
}

// Metric is the type for training metric (e.g. acc).
type Metric struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Hyperparameter is the type for training hyperparameter (e.g. learning rate).
type Hyperparameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Signature struct {
	Inputs  []Tensor `json:"inputs,omitempty" yaml:"inputs,omitempty"`
	Outputs []Tensor `json:"outputs,omitempty" yaml:"outputs,omitempty"`
	Layers  map[string]int  `json:"layers,omitempty" yaml:"layers,omitempty"`
}

type Tensor struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	Size  []int  `json:"size,omitempty" yaml:"size,omitempty"`
	DType string `json:"dtype,omitempty" yaml:"dtype,omitempty"`
	// OpType is special for PMML
	OpType string `json:"opType,omitempty" yaml:"opType,omitempty"`
	// Values is special for PMML
	Values []string `json:"values,omitempty" yaml:"values,omitempty"`
}

type Training struct {
	Git GitRepo `json:"git,omitempty" yaml:"git,omitempty"`
}

type Dataset struct {
	Git GitRepo `json:"git,omitempty" yaml:"git,omitempty"`
}

type GitRepo struct {
	Repository string `json:"repository,omitempty" yaml:"repository,omitempty"`
	Revision   string `json:"revision,omitempty" yaml:"revision,omitempty"`
}
