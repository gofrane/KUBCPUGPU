package main

import (
	"errors"
	"log"
)




/*
type InformationTableType struct {
	Id int
	Image string
	NodeUsed Node
	timeImagePerNode float64


}
*/








//*************************//
//// classify the node depeond on the resource type
func getNodeType() (nodeCPU []Node, nodeGPU []Node, err error) {
	nodes, errNodes := getNodes()

	if errNodes != nil {
		err := errors.New("no node available ")
		return nil, nil, err

	} else {

		for _, n := range nodes.Items {

			s := n.Status.Capacity

			for k, _ := range s {
				if _, ok := s["nvidia.com/gpu"]; ok {
					nodeGPU = append(nodeGPU, n)
					_ = k
					break
				} else {
					nodeCPU = append(nodeCPU, n)
					break

				}

			}
		}

	}

	return nodeGPU, nodeCPU, nil

}

///////////////////////////////////////////////////

//check the node type of the node elected
func checkNodeType(GPUnode []Node, bestNode string) string {

	check := false
	for _, node := range GPUnode {

		if node.Metadata.Name == bestNode {

			check = true

			break

		}

	}
	if check {

		return "GPU"
	} else {

		return "CPU"
	}

}
/////////////////////////////////////////////////





//****************************************************//
func decison(nodes []Node, pod *Pod, ImagePending string,ImageCPU string,ImageGPU string,InformationTable []dataBaseInformationTable) (*Pod,Node, error) {
	//make table node type
	var err error
	var bestNode Node
	var newPod *Pod
	GPUnode, _, getNodeTyperror:= getNodeType() // get the node type
// if getNodeTyperror not nil
	if getNodeTyperror !=nil {

		log.Println(getNodeTyperror)
		return nil,Node{},getNodeTyperror

	}


	if (len(InformationTable)>0) {


		bestNode = electBestNode(nodes,ImagePending,InformationTable,GPUnode ,pod ) //elect the best node

		typeNode := checkNodeType(GPUnode, bestNode.Metadata.Name) //check the best node seleted after the best node


	 //check the type of node selected

         	if typeNode == "GPU" {
            //modify the pod for the GPU
		   newPod=modifyPodforGPU(pod,ImageGPU)

	      }else{

				//modify the pod for the CPU
	      	newPod=modifyPodforCPU(pod,ImageCPU)

	}
	        return newPod,bestNode, nil
      }else{
	   err=errors.New("the decison can not work")
	   return newPod,Node{},err
       }


}
