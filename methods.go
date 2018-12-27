package main
/*
type PodPhaseList struct {
	PodRunning []Pod
	PodSucceeded []Pod
	PodPending []Pod
	PodFailed []Pod

}
*/
func modifyPodforCPU(oldPod *Pod, ImagePendingSelected string) *Pod {

	//var newPod Pod

	oldContainer := oldPod.Spec.Containers[0]

	newPod := &Pod{
		Kind:     oldPod.Kind,
		Metadata: oldPod.Metadata,
		Spec: PodSpec{
			NodeName: oldPod.Spec.NodeName,
			Containers: []Container{
				{
					Name: oldContainer.Name,
					Resources: ResourceRequirements{

						Limits:   map[string]string{"nvidia.com/gpu": "1"},
						Requests: oldContainer.Resources.Requests,
					},

					Image: ImagePendingSelected,
				},
			},
		},
		Phase:  oldPod.Phase,
		Status: oldPod.Status,
	}
	return newPod
}



func modifyPodforGPU(oldPod *Pod, ImagePendingSelected string) *Pod {

	//var newPod Pod

	oldContainer := oldPod.Spec.Containers[0]

	newPod := &Pod{
		Kind:     oldPod.Kind,
		Metadata: oldPod.Metadata,
		Spec: PodSpec{
			NodeName: oldPod.Spec.NodeName,
			Containers: []Container{
				{
					Name: oldContainer.Name,
					Resources: ResourceRequirements{

						Limits:   map[string]string{"nvidia.com/gpu": "1"},
						Requests: oldContainer.Resources.Requests,
					},

					Image: ImagePendingSelected,
				},
			},
		},
		Phase:  oldPod.Phase,
		Status: oldPod.Status,
	}
	return newPod
}
/*
https://golangr.com/difference-between-two-dates/

http://www.golangprograms.com/get-hours-days-minutes-and-seconds-difference-between-two-dates-future-and-past.html

https://stackoverflow.com/questions/40260599/difference-between-two-time-time-objects/40260666
*/



/*
func (PodPhaseList *PodPhaseList)PodListClassify(lastPod Pod) error{


	var err error

	phase := lastPod.Status.Phase

	switch phase {

	case "Running":
		PodPhaseList.PodRunning = append(PodPhaseList.PodRunning, lastPod)

	case "Pending":
		PodPhaseList.PodPending = append(PodPhaseList.PodPending, lastPod)

	case "Succeeded":
		PodPhaseList.PodSucceeded = append(PodPhaseList.PodSucceeded, lastPod)

	case "Failed":
		PodPhaseList.PodFailed = append(PodPhaseList.PodFailed, lastPod)

	default:
		err=errors.New("Phase  does not exist")

	}

	fmt.Println(PodPhaseList)

	return err
}

*/