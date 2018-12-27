package main

import "time"

// Event is a report of an event somewhere in the cluster.
type Event struct {
	ApiVersion     string          `json:"apiVersion,omitempty"`
	Count          int64           `json:"count,omitempty"`
	FirstTimestamp string          `json:"firstTimestamp"`
	LastTimestamp  string          `json:"lastTimestamp"`
	InvolvedObject ObjectReference `json:"involvedObject"`
	Kind           string          `json:"kind,omitempty"`
	Message        string          `json:"message,omitempty"`
	Metadata       Metadata        `json:"metadata"`
	Reason         string          `json:"reason,omitempty"`
	Source         EventSource     `json:"source,omitempty"`
	Type           string          `json:"type,omitempty"`
}

// EventSource contains information for an event.
type EventSource struct {
	Component string `json:"component,omitempty"`
	Host      string `json:"host,omitempty"`
}

// ObjectReference contains enough information to let you inspect or modify
// the referred object.
type ObjectReference struct {
	ApiVersion string `json:"apiVersion,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Name       string `json:"name,omitempty"`
	Namespace  string `json:"namespace,omitempty"`
	Uid        string `json:"uid"`
}

// PodList is a list of Pods.
type PodList struct {
	ApiVersion string       `json:"apiVersion"`
	Kind       string       `json:"kind"`
	Metadata   ListMetadata `json:"metadata"`
	Items      []Pod        `json:"items"`
}

type PodWatchEvent struct {
	Type   string `json:"type"`
	Object Pod    `json:"object"`
}

type Pod struct {
	ApiVersion string   `json:"apiVersion ,omitempty"`
	Kind       string   `json:"kind,omitempty"`
	Metadata   Metadata `json:"metadata"`
	Spec       PodSpec  `json:"spec"`
	Phase      string   `json:"phase"`
	Status     Status   `json:"status"`
}

type Status struct {
	Phase             string            `json:"phase"`
	Conditions        []Conditions      `json:"conditions"`
	HostIP            string            `json:"hostIP"`
	PodIP             string            `json:"podIP"`
	StartTime         *time.Time        `json:"startTime"`
	CompletionTime    *time.Time        `json:"completionTime"`
	ContainerStatuses []ContainerStatus `json:"containerStatuses"`
	LastState         []ContainerState  `json:"lastState"`
	QosClass          string            `json:"qosClass"`
}

type ContainerState struct {
	Running    ContainerStateRunning    `json:"running"`
	Terminated ContainerStateTerminated `json:"terminated" `
	Waiting    ContainerStateWaiting    `json:"waiting" `
}

type ContainerStateRunning struct {
	StartedAt *time.Time `json:"startedAt"`
}

type ContainerStateTerminated struct {
	FinishedAt *time.Time `json:"finishedAt"`
	StartedAt  *time.Time `json:"startedAt"`
}

type ContainerStateWaiting struct {
	Message string `json:"message"`
	Reason  string `json:"reason"`
}

type Conditions struct {
	Type               string      `json:"type"`
	Status             string      `json:"status"`
	LastProbeTime      interface{} `json:"lastProbeTime"`
	LastTransitionTime *time.Time  `json:"lastTransitionTime"`
}

type ContainerStatus struct {
	Name  string `json:"name"`
	State struct {
		Running struct {
			StartedAt *time.Time `json:"startedAt"`
		} `json:"running"`
	} `json:"state"`
	LastState struct {
	} `json:"lastState"`
	Ready        bool   `json:"ready"`
	RestartCount int    `json:"restartCount"`
	Image        string `json:"image"`
	ImageID      string `json:"imageID"`
	ContainerID  string `json:"containerID"`
}

type PodSpec struct {
	NodeName   string      `json:"nodeName"`
	Containers []Container `json:"containers"`
}

type Container struct {
	Name      string               `json:"name"`
	Resources ResourceRequirements `json:"resources"`
	Image     string               `json:"image"`
}

type ResourceRequirements struct {
	Limits   ResourceList `json:"limits"`
	Requests ResourceList `json:"requests"`
}

type ResourceList map[string]string

type Binding struct {
	ApiVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Target     Target   `json:"target"`
	Metadata   Metadata `json:"metadata"`
}

type Target struct {
	ApiVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Name       string `json:"name"`
}

type NodeList struct {
	ApiVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Items      []Node
}

type Node struct {
	Metadata Metadata   `json:"metadata"`
	Status   NodeStatus `json:"status"`
}

type NodeStatus struct {
	Capacity    ResourceList `json:"capacity"`
	Allocatable ResourceList `json:"allocatable"`
}

type ListMetadata struct {
	ResourceVersion string `json:"resourceVersion"`
}

type Metadata struct {
	Name              string            `json:"name"`
	GenerateName      string            `json:"generateName"`
	ResourceVersion   string            `json:"resourceVersion"`
	Labels            map[string]string `json:"labels"`
	Annotations       map[string]string `json:"annotations"`
	Uid               string            `json:"uid"`
	CreationTimestamp *time.Time        `json:"creationTimestamp"`
}
