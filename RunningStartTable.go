package main

import "time"

func appendRunningPodTable(RunningPodslist *PodList, initial int , final int)[]time.Time{

	var  RunningPod []time.Time


	for i := initial; i < final; i++ {

		RunningPod = append(RunningPod, *RunningPodslist.Items[i].Status.StartTime)

	}

	return RunningPod

}

func CheckRunningStartTable ( currentPodTime time.Time,RunningPodTime []time.Time)bool{


	verif:=false
	for _,v:=range RunningPodTime {

		if v.Sub(currentPodTime)==0 {


			verif=true
			break
		}

	}
	return  verif
}
