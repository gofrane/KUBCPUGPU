package main
type sortListPod []Pod
func (a sortListPod) Len() int           { return len(a) }
func (a sortListPod) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortListPod) Less(i, j int) bool {
 t1:=*a[i].Status.StartTime
 t2:=*a[j].Status.StartTime
	diff := t2.Sub(t1)

	duration := float64(diff.Seconds())
if duration > 0{
	return true
}
	return false
}



/*
https://stackoverflow.com/questions/28999735/what-is-the-shortest-way-to-simply-sort-an-array-of-structs-by-arbitrary-field
 */