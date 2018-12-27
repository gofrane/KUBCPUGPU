package main

import "time"

//store the waiting pod
//////////////////////////////////////
type WaitingPod struct {

	podName string
	nodeSelected string
	durationEstimated float64
	waitingTime float64
	StarttimeEstimated *time.Time
	EstimatedtimeTermination *time.Time

}

type WaitingPodList struct {
	List []WaitingPod
}


var WaitingGPUPod WaitingPodList


//////////////////////////////////////////////



func(WaitingGPUPod  *WaitingPodList) addToWaitingList(pod string, node string ,totalTime float64 ) {


	var  waitingPod  WaitingPod

	waitingPod.podName=pod
	waitingPod.nodeSelected=node
	waitingPod.waitingTime=totalTime


	WaitingGPUPod.List  =append(WaitingGPUPod.List ,waitingPod )


}


func ( WaitingGPUPod  WaitingPodList)totalWaitaingTimeGPUnode(node string)float64{

var toltalWaitingTime float64


	for _,v:= range WaitingGPUPod.List{
		if v.nodeSelected==node {
			toltalWaitingTime += v.waitingTime
		}

	}

return  toltalWaitingTime

}




func (WaitingGPUPod  *WaitingPodList)deleteWaitingGPUPod(SucceededPods *PodList  ,RunningList CurrentRunningList) {


	if len(RunningList.List)>0 {

		for k,v:=range WaitingGPUPod.List {

			for _,x:=range RunningList.List{

				if v.podName==x.PodName {
					if len( WaitingGPUPod.List) >1 {
						WaitingGPUPod.List = append( WaitingGPUPod.List[:k],  WaitingGPUPod.List[k+1:]...)
					}else{
						WaitingGPUPod.List=nil
					}
				}
			}

		}


	}
	for k,v:=range WaitingGPUPod.List {

		for _,x:=range SucceededPods.Items{

			if v.podName==x.Metadata.Name {
				if len( WaitingGPUPod.List) >1 {
					WaitingGPUPod.List = append( WaitingGPUPod.List[:k],  WaitingGPUPod.List[k+1:]...)
				}else{
					WaitingGPUPod.List=nil
				}
			}
		}

	}

}