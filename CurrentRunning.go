package main

import (
	"fmt"
	"time"
)

type CurrentRunning struct {
	PodName            string
	NodeName           string
	timeStartPod       time.Time
	timeStartRunning   time.Time
	durartionEstimated float64
	FinalTimeEstimated time.Time
}

type CurrentRunningList struct {
	List []CurrentRunning
}

var RunningList CurrentRunningList

/////////////////////////////////////////////
// this function append the  CurrentRunningList which has the information of the ruuning pod
func (RunningList *CurrentRunningList) appendCurrentRunning(lastPod Pod, timeStartRunning time.Time) {

	durationEstimated := getDurationEstimated(lastPod) // last-k-1 used to modify theo order of the table (the last ruuning pod on the last index of the table)

	FinalTimeEstimated := getFinalTimeEstimated(lastPod, durationEstimated)
	CurrentRunning := CurrentRunning{PodName: lastPod.Metadata.Name,
		NodeName:           lastPod.Spec.NodeName,
		timeStartPod:       *lastPod.Status.StartTime,
		timeStartRunning:   timeStartRunning,
		durartionEstimated: durationEstimated,
		FinalTimeEstimated: FinalTimeEstimated}

	RunningList.List = append(RunningList.List, CurrentRunning)

}

////////////////////////////////////////////

///
func (RunningList *CurrentRunningList) deletePodCurrentRunningList(SucceededPods *PodList) {
	//SucceededPods,_:=getSucceededPods()
	//FailedPods,_:=getFailedPods()
	for k, v := range RunningList.List {

		for _, x := range SucceededPods.Items {
			fmt.Println(" len (RunningList.List ) ", len(RunningList.List))
			if v.PodName == x.Metadata.Name {
				if len(RunningList.List) > 1 {
					RunningList.List = append(RunningList.List[:k], RunningList.List[k+1:]...)
				} else {
					RunningList.List = nil
				}
			}
		}
		/*
			for _,y:=range FailedPods.Items{

				if v.PodName==y.Metadata.Name {

					CurrentRunningList=append(CurrentRunningList[:k],CurrentRunningList[:k+1]...)


				}
		*/
	}

}

////////////////////////////////////////////

func (RunningList CurrentRunningList) RemainingTimeGPUnode(node string) float64 {

	var RemainingTime float64

	timeNow := time.Now().UTC()
	for _, v := range RunningList.List {

		if v.NodeName == node {
			diff := (v.FinalTimeEstimated).Sub(timeNow)
			RemainingTime = float64(diff.Seconds())

			//MdurationEstimated := time.Duration(durationRemained)

		}
	}

	return RemainingTime

}
