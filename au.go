package main

/*
type CurrentRunning struct {

	PodName string
	NodeName string
	timeStartPod *time.Time
	timeStartRunning time.Time
	durartionEstimated float64
	FinalTimeEstimated  time.Time

}

var (
	processorLock       = &sync.Mutex{}
	succeeddedPodLength int
	failedPodLength     int
	runningPodLength     int
	verifSucceedded     = true
	verifFailed         = true
	verifRunning        = true
	RunningStartTable      []time.Time
	timeStartRunning        time.Time
	CurrentRunningList []CurrentRunning
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
					runningPodLength= len(RunningPods.Items)
					addPodtoDBrunning(0, runningPodLength, RunningPods,timeStartRunning)
					RunningStartTable=appendRunningPodTable(RunningPods, 0 , len(RunningPods.Items))
					appendCurrentRunning(CurrentRunningList,RunningPods,0,len(RunningPods.Items))
					verifRunning = false


				}


				if runningPodLength != len(RunningPods.Items) {
					timeStartRunning= time.Now().UTC()
					if CheckRunningStartTable(*RunningPods.Items[0].Status.StartTime,RunningStartTable) ==false{
						appendCurrentRunning(CurrentRunningList,RunningPods,len(CurrentRunningList),len(CurrentRunningList)+1)

						// RunningStartTable = appendRunningPodTable(RunningPods, len(RunningPods.Items), len(RunningPods.Items)+1)
						RunningStartTable = append(RunningStartTable, *RunningPods.Items[0].Status.StartTime)
						addPodtoDBrunning(0, 1, RunningPods,timeStartRunning)

						//runningPodLength  take new length of RunningPods
						//runningPodLength = len(RunningPods.Items)

					}

				}


				///////////////////////////////////////////
				SucceededPods, errSucceededPod := getSucceededPods()
				if errSucceededPod == nil {
					if verifSucceedded {
						succeeddedPodLength = len(SucceededPods.Items)
						addPodtoDB(0, succeeddedPodLength, SucceededPods)
						verifSucceedded = false
					}

					if succeeddedPodLength < len(SucceededPods.Items) {

						addPodtoDB(succeeddedPodLength, len(SucceededPods.Items), SucceededPods)

						succeeddedPodLength = len(SucceededPods.Items)

					}
				}
				/////////////////////////////////

				deletePodterminated (CurrentRunningList  )
				fmt.Println(CurrentRunningList )
				///////////////////
			}
		}
	}
}

/////////////////////////////////////////////
// this function append the  CurrentRunningList which has the information of the ruuning pod
func appendCurrentRunning ( CurrentRunningList []CurrentRunning ,RunningPods *PodList,first int , last int ){

	for k:=first ;k<last ;k++{
		durationEstimated:=getDuration(RunningPods.Items[last-k-1]) // last-k-1 used to moddify theo order of the table (the last ruuning pod on the last index of the table)

		FinalTimeEstimated:=getFinalTimeEstimated(RunningPods.Items[last-k-1] ,durationEstimated )
		CurrentRunningList[k].PodName=RunningPods.Items[last-k-1].Metadata.Name
		CurrentRunningList[k].NodeName=RunningPods.Items[last-k-1].Spec.NodeName
		CurrentRunningList[k].timeStartPod=RunningPods.Items[last-k-1].Metadata.CreationTimestamp
		CurrentRunningList[k].durartionEstimated=durationEstimated
		CurrentRunningList[k].FinalTimeEstimated=FinalTimeEstimated


	}
	fmt.Println(CurrentRunningList)

}

////////////////////////////////////////////

///
func deletePodterminated (CurrentRunningList []CurrentRunning ) {
	SucceededPods,_:=getSucceededPods()
	FailedPods,_:=getFailedPods()
	for k,v:=range CurrentRunningList {

		for _,x:=range SucceededPods.Items{

			if v.PodName==x.Metadata.Name {
				CurrentRunningList=append(CurrentRunningList[:k],CurrentRunningList[:k+1]...)
			}
		}

		for _,y:=range FailedPods.Items{

			if v.PodName==y.Metadata.Name {

				CurrentRunningList=append(CurrentRunningList[:k],CurrentRunningList[:k+1]...)


			}
		}

	}



}



////////////////////////////////////////////
func appendRunningPodTable(RunningPodslist *PodList, initial int , final int)[]time.Time{

	var  RunningPod []time.Time


	for i := initial; i < final; i++ {

		RunningPod = append(RunningPod, *RunningPodslist.Items[i].Status.StartTime)

	}

	return RunningPod

}

func CheckRunningStartTable ( currentPodTime time.Time,RunningPodTime []time.Time)bool{

	fmt.Println("RunningPodTime "   , len(RunningPodTime))
	verif:=false
	for _,v:=range RunningPodTime {

		if v.Sub(currentPodTime)==0 {

			verif=true
			break
		}

	}
	return  verif
}











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

	nodes, err := fit(pod)
	if err != nil {
		return err
	}
	if len(nodes) == 0 {
		return fmt.Errorf("Unable to schedule pod (%s) failed to fit in any node", pod.Metadata.Name)
	}

	imagePending := getImage(pod)
	imagePendingCPU, ImagePendingGPU, _ := db.imageSearch(imagePending)

	//	var InformationTable []InformationTableType
	//row1 := InformationTableType{0, "gofrane/cudatest", nodes[0],  10.0}
	//row2 := InformationTableType{1, "gofrane/cudatest", nodes[1],  7.0}

	//InformationTable = append(InformationTable, row1)
	//InformationTable = append(InformationTable, row2)
	InformationTable := db.informationSearch(imagePending)
	//fmt.Println(InformationTable)

	pod, node, err := decison(nodes, pod, imagePending, imagePendingCPU, ImagePendingGPU, InformationTable)
	//fmt.Println(pod.Spec.Containers[0].Resources.Limits)

	if err != nil {
		return err
	}

	// pod.Spec.Containers[0].Resources.Limits= map[string]string{"nvidia.com/gpu":"1"}
	err = bind(pod, node)
	if err != nil {
		return err
	}
	return nil
}

func schedulePods() error {
	processorLock.Lock()
	defer processorLock.Unlock()

	///////////////////////////////////////////
	SucceededPods, errSucceededPod := getSucceededPods()
	if errSucceededPod == nil {
		if verifSucceedded {
			succeeddedPodLength = len(SucceededPods.Items)
			addPodtoDB(0, succeeddedPodLength, SucceededPods)
			verifSucceedded = false
		}

		if succeeddedPodLength < len(SucceededPods.Items) {

			addPodtoDB(succeeddedPodLength, len(SucceededPods.Items), SucceededPods)

			succeeddedPodLength = len(SucceededPods.Items)

		}
	}
	/////////////////////////////////

	FailedPods, errFailedPod := getFailedPods()
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

	////////////////////////////


	//










	////////////////////////////////

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
*/