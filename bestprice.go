package main

//func SwitchDeviceType(newPod *Pod){

//newPod.Spec.Containers[0].Resources.Limits= map[string]string{"nvidia.com/gpu":"1"}

// }

/* func bestPrice(nodes []Node) (Node, error) {
	type NodePrice struct {
		Node  Node
		Price float64
	}
	//SwitchDeviceType(pod)


	var bestNodePrice *NodePrice
	for _, n := range nodes {
		price, ok := n.Metadata.Annotations["hightower.com/cost"]
		if !ok {
			continue
		}
		f, err := strconv.ParseFloat(price, 32)
		if err != nil {
			return Node{}, err
		}
		if bestNodePrice == nil {
			bestNodePrice = &NodePrice{n, f}
			continue
		}
		if f < bestNodePrice.Price {
			bestNodePrice.Node = n
			bestNodePrice.Price = f
		}
	}

	if bestNodePrice == nil {
		bestNodePrice = &NodePrice{nodes[0], 0}
	}

	return bestNodePrice.Node, nil
}
*/
