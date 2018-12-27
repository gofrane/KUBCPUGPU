package main

import "time"

type Image struct {
	MainImage string `json:"MainImageName"`
	CPUimage  string `json:"CPUImageName"`
	GPUimage  string `json:"GPUImageName"`
}

type dataBaseInformationTable struct {
	Id          int     `json:"id"`
	ImageUsed   string  `json:"image"`
	NodeUsed    string  `json:"node"`
	RunningTime float64 `json:"timeEstimated"`
}

type Statistic struct {
	Id             int        `json:"id"`
	ImageUsed      string     `json:"imageUsed"`
	NodeUsed       string     `json:"nodeUsed"`
	PodFinalStatus *time.Time `json:"podFinalStatus"`
	PendingTime *time.Time `json:"pendingTime"`
	StartDeployementTime *time.Time ` json:"startDeployementTime"`
	FinalDeployementTime *time.Time `json:"finalDeployementTime"`
	//PendingDuration float64 `json:"pendingtDuration"`
	duration float64 `json:"TotalDuration"`
}
