package main

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"time"
)

/*
type CurrentRunning struct {

	PodName string
	NodeName string
	timeStartPod time.Time
	timeStartRunning time.Time
	durartionEstimated float64
	FinalTimeEstimated  time.Time

}

type CurrentRunningList struct {
	List []CurrentRunning
}
*/
var (
	processorLock          = &sync.Mutex{}
	succeeddedPodLength    int
	failedPodLength        int
	runningPodLength       int
	verifSucceedded        = true
	verifFailed            = true
	verifRunning           = true
	RunningStartTable      []time.Time
	timePendingTable       []time.Time
	timeFinishPendingTable []time.Time
	timePodMatchedTable    []time.Time
	timePodMatched         time.Time
	//RunningList CurrentRunningList
	//RunningPodTableLength int
)

func reconcileUnscheduledPods(interval int, done chan struct{}, wg *sync.WaitGroup) {
	for {
		select {
		case <-time.After(time.Duration(interval) * time.Second):
			err := schedulePods()
			if err != nil {
				log.Println(err)
			}
		case <-done:
			wg.Done()
			log.Println("Stopped reconciliation loop.")
			return
		default:
			RunningPods, errRunningPod := getRunningPods()

			if errRunningPod == nil {
				if verifRunning {
					runningPodLength = len(RunningPods.Items)
					//addPodtoDBrunning(0, 1, RunningPods,time.Time{})

					RunningStartTable = appendRunningPodTable(RunningPods, 0, len(RunningPods.Items)) //the RunningStartTable has the start tome of each pod

					verifRunning = false

				}

				if runningPodLength != len(RunningPods.Items) { // if the ruuning PodLength change that means add or remobe 1 pod for the running lsit
					timeStartRunning := time.Now().UTC() // approximately it is the time of start running after container creation

					sort.Sort(sort.Reverse(sortListPod(RunningPods.Items))) //Sort the RunningPods
					if CheckRunningStartTable(*RunningPods.Items[0].Status.StartTime, RunningStartTable) == false {
						RunningList.appendCurrentRunning(RunningPods.Items[0], timeStartRunning)
						// RunningStartTable = appendRunningPodTable(RunningPods, len(RunningPods.Items), len(RunningPods.Items)+1)
						RunningStartTable = append(RunningStartTable, *RunningPods.Items[0].Status.StartTime)

						addPodtoDBrunning(0, 1, RunningPods, timeStartRunning)

						//runningPodLength  take new length of RunningPods
						//runningPodLength = len(RunningPods.Items)

					}

				}

			}
		}
	}
}

/*
func ( RunningList *CurrentRunningList )appendCurrentRunning(lastPod Pod ,timeStartRunning time.Time ){



		durationEstimated:=getDuration(lastPod) // last-k-1 used to modify theo order of the table (the last ruuning pod on the last index of the table)

		FinalTimeEstimated:=getFinalTimeEstimated(lastPod ,durationEstimated )
	CurrentRunning:=CurrentRunning{PodName:lastPod.Metadata.Name,
	                               NodeName:lastPod.Spec.NodeName,
	                               timeStartPod: *lastPod.Status.StartTime,
		                           timeStartRunning:timeStartRunning,
	                               durartionEstimated:durationEstimated,
	                               FinalTimeEstimated:FinalTimeEstimated}

                      RunningList.List=append(RunningList.List,CurrentRunning)

}

////////////////////////////////////////////

///
func ( RunningList *CurrentRunningList )deletePodCurrentRunningList(SucceededPods *PodList   ) {
	//SucceededPods,_:=getSucceededPods()
	//FailedPods,_:=getFailedPods()
	for k,v:=range RunningList.List {

		for _,x:=range SucceededPods.Items{

			if v.PodName==x.Metadata.Name {
               if len(RunningList.List) >1 {
				   RunningList.List = append(RunningList.List[:k], RunningList.List[k+1:]...)
			   }else{
			   	 RunningList.List=nil
			   }
			}
		}
/*
		for _,y:=range FailedPods.Items{

			if v.PodName==y.Metadata.Name {

				CurrentRunningList=append(CurrentRunningList[:k],CurrentRunningList[:k+1]...)


			}

		}

	}


*/

////////////////////////////////////////////

//////////////////////////////////////////////////

func monitorUnscheduledPods(done chan struct{}, wg *sync.WaitGroup) {
	pods, errc := watchUnscheduledPods()

	for {
		select {
		case err := <-errc:
			log.Println(err)
		case pod := <-pods:
			processorLock.Lock()
			time.Sleep(2 * time.Second)
			err := schedulePod(&pod)
			if err != nil {
				log.Println(err)
			}
			processorLock.Unlock()
		case <-done:
			wg.Done()
			log.Println("Stopped scheduler.")
			return
		}
	}
}

func schedulePod(pod *Pod) error {
	var mutex sync.Mutex

	mutex.Lock()
	defer mutex.Unlock()
	newTimeStartPending := time.Now().UTC()
	timePendingTable = append(timePendingTable, newTimeStartPending)

	fmt.Println("time Start Pending  ", newTimeStartPending)
	SucceededPods, errSucceededPod := getSucceededPods()
	if (len(RunningList.List) > 0) && (len(SucceededPods.Items) > 0) {
		RunningList.deletePodCurrentRunningList(SucceededPods)

	}
	if errSucceededPod == nil {
		if verifSucceedded {
			succeeddedPodLength = len(SucceededPods.Items)
			//	PendingLength:=len(timePendingTable)
			//	addPodtoDB(0, succeeddedPodLength, SucceededPods,timePendingTable[PendingLength-1],timePodMatched )
			verifSucceedded = false
		}

		if succeeddedPodLength < len(SucceededPods.Items) {
			PendingLength := len(timePendingTable)
			sort.Sort(sort.Reverse(sortListPod(SucceededPods.Items)))
			fmt.Println(SucceededPods.Items[0])
			addPodtoDB(0, len(SucceededPods.Items)-succeeddedPodLength, SucceededPods, timePendingTable[PendingLength-1], timePodMatchedTable[PendingLength-1])

			succeeddedPodLength = len(SucceededPods.Items)

		}
	}
	nodes, err := fit(pod)
	if err != nil {
		return err
	}
	if len(nodes) == 0 {
		return fmt.Errorf("Unable to schedule pod (%s) failed to fit in any node", pod.Metadata.Name)
	}

	if len(WaitingGPUPod.List) > 0 {
		WaitingGPUPod.deleteWaitingGPUPod(SucceededPods, RunningList)
	}

	imagePending := getImage(pod)
	imagePendingCPU, ImagePendingGPU, _ := db.imageSearch(imagePending)

	InformationTable := db.informationSearch(imagePending)

	pod, node, err := decison(nodes, pod, imagePending, imagePendingCPU, ImagePendingGPU, InformationTable)

	config.UpdatePod2(pod.Metadata.Name)
	//fmt.Println("the newwwwwwwwwwww",pods)
	//fmt.Println(pods)
	//if err != nil {
	//	return err
	//}

	err = bind(pod, node)
	timePodMatched = time.Now().UTC()
	timePodMatchedTable = append(timePodMatchedTable, timePodMatched)

	if err != nil {
		return err
	}
	return nil
}

func schedulePods() error {
	processorLock.Lock()
	defer processorLock.Unlock()

	////////////////////////
	/*	FailedPods, errFailedPod := getFailedPods()
		if errFailedPod == nil {
			if verifFailed {
				failedPodLength = len(FailedPods.Items)
				addPodtoDB(0, failedPodLength, FailedPods)
				verifFailed = false
			}

			if failedPodLength < len(FailedPods.Items) {

				addPodtoDB(failedPodLength, len(FailedPods.Items), FailedPods)

				failedPodLength = len(FailedPods.Items)

			}

		}


	*/

	pods, err := getUnscheduledPods()
	if err != nil {
		return err
	}

	for _, pod := range pods {
		err := schedulePod(pod)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}
