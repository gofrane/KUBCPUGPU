package main

import "fmt"







func availableNodeperImage(nodes []Node ,nodeCheck string)bool{


	verif :=false
	for _,v:=range nodes {

		if nodeCheck == v.Metadata.Name{
			verif=true
			break

		}

	}
	return  verif
}

//***************************************************//

//create image per node table (ImagePending +Node avialable )
func ImageperNodeAvailable (ImagePending string, InformationTable  []dataBaseInformationTable , nodes []Node)[]dataBaseInformationTable{

	var ImageperNodeTable []dataBaseInformationTable

	for _,v:=range InformationTable {

		if v.ImageUsed==ImagePending && availableNodeperImage(nodes,v.NodeUsed ){

			ImageperNodeTable = append(ImageperNodeTable, v)
		}


	}
	fmt.Println(ImageperNodeTable)
	return  ImageperNodeTable

}







//*********************************************************//

func electBestNode(nodes []Node, ImagePending string, InformationTable []dataBaseInformationTable ,GPUnode []Node,pod *Pod) Node {
	var totalTimeForWaiting float64
	var Totaltime float64
	var RemainingTime float64
	var timeWaitingGPUPod float64
	var bestNode Node

	ImageperNodeTable:=ImageperNodeAvailable(ImagePending,InformationTable,nodes) //the slice of the image per node available
	fmt.Println("Image per node table",ImageperNodeTable )

	index:=0
	BestTime:=ImageperNodeTable[0].RunningTime
	i:=0






	///////////////////////////////
	for(i<len(ImageperNodeTable)-1) {

		if (len(WaitingGPUPod.List) == 0 && len(RunningList.List) == 0) {
			if BestTime > ImageperNodeTable[i+1].RunningTime {
				BestTime = ImageperNodeTable[i+1].RunningTime
				index = i + 1
               fmt.Println("Best time 1", BestTime)
				totalTimeForWaiting=0
			}
			i = i + 1
			continue
		} else if len(WaitingGPUPod.List) == 0 && len(RunningList.List) > 0 {
			if checkNodeType(GPUnode, ImageperNodeTable[i+1].NodeUsed) == "GPU" {
				RemainingTime = RunningList.RemainingTimeGPUnode(ImageperNodeTable[i+1].NodeUsed)
			} else {
				RemainingTime = 0
			}
			Totaltime = RemainingTime + ImageperNodeTable[i+1].RunningTime
			if BestTime > Totaltime {
				BestTime = Totaltime
				index = i + 1
				fmt.Println("Best time 2", BestTime)

				totalTimeForWaiting=Totaltime

			}
			i = i + 1
			continue

		} else if len(WaitingGPUPod.List) > 0 && len(RunningList.List) == 0 {
			if (checkNodeType(GPUnode, ImageperNodeTable[i+1].NodeUsed)) == "GPU" {
				timeWaitingGPUPod = WaitingGPUPod.totalWaitaingTimeGPUnode(ImageperNodeTable[i+1].NodeUsed)
			} else {
				timeWaitingGPUPod = 0
			}
			Totaltime = timeWaitingGPUPod + ImageperNodeTable[i+1].RunningTime
			if BestTime > Totaltime {
				BestTime = Totaltime
				index = i + 1
				totalTimeForWaiting=Totaltime
				fmt.Println("Best time 2", BestTime)

			}
			i = i + 1
			continue

		} else if len(WaitingGPUPod.List) > 0 && len(RunningList.List) > 0 {

			if (checkNodeType(GPUnode, ImageperNodeTable[i+1].NodeUsed)) == "GPU" {
				RemainingTime = RunningList.RemainingTimeGPUnode(ImageperNodeTable[i+1].NodeUsed)
				timeWaitingGPUPod = WaitingGPUPod.totalWaitaingTimeGPUnode(ImageperNodeTable[i+1].NodeUsed)

			} else {
				timeWaitingGPUPod = 0
				RemainingTime = 0

			}
			Totaltime = timeWaitingGPUPod + RemainingTime+ImageperNodeTable[i+1].RunningTime
			if BestTime > Totaltime {
				BestTime = Totaltime
				index = i + 1
				totalTimeForWaiting=Totaltime
				fmt.Println("Best time 3 ", BestTime)


			}
			i = i + 1


		}

	}




	bestNodeName:= ImageperNodeTable[index].NodeUsed
fmt.Println("bestNodeName ",	bestNodeName )
	fmt.Println("totalTimeForWaiting ",	totalTimeForWaiting )
	 if checkNodeType(GPUnode, bestNodeName) == "GPU" {
		 if totalTimeForWaiting > 0 {

			 WaitingGPUPod.addToWaitingList(pod.Metadata.Name, bestNodeName, totalTimeForWaiting)
		 }
	 }
	//match the name of node selected to the node type
	for _,n:=range  nodes{
		if bestNodeName==n.Metadata.Name{
			bestNode =n
			break

		}

	}
	return bestNode

}





