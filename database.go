package main

import (
	"log"
	"time"
)

/*
type dbConn struct {
	db *sql.DB
}
*/
/*
func (db dbConn)openDataBase()error{

var err error

	db.db, err = sql.Open("mysql", "gofrane:gofrane@tcp(127.0.0.1:3306)/test")
	if err != nil {
		log.Fatal(err)
	}

	 return err

}
*/

func (db dbConn) imageSearch(imagePending string) (string, string, error) {

	results, err := db.db.Query("SELECT MainImageName , CPUimageName ,  GPUimageName FROM imagetable where MainImageName=? ", imagePending)
	if err != nil {
		panic(err.Error())
	}
	var imCPU, imGPU string

	for results.Next() {
		var image Image
		// for each row, scan the result into our image composite object
		err = results.Scan(&image.MainImage, &image.CPUimage, &image.GPUimage)
		if err != nil {
			panic(err.Error())
		}

		log.Printf(image.MainImage)
		log.Printf(image.CPUimage)
		log.Printf(image.GPUimage)
		imCPU = image.CPUimage
		imGPU = image.GPUimage
	}
	defer results.Close()

	return imCPU, imGPU, err
}

func (db dbConn) informationSearch(imagePending string) []dataBaseInformationTable {

	results, err := db.db.Query("SELECT id , ImageName ,  NodeName , DeployementTime  FROM informationtable where ImageName=? ", imagePending)
	if err != nil {
		panic(err.Error())
	}
	var InfomationTable []dataBaseInformationTable

	for results.Next() {
		var information dataBaseInformationTable
		// for each row, scan the result into our image composite object
		err := results.Scan(&information.Id, &information.ImageUsed, &information.NodeUsed, &information.RunningTime)
		if err != nil {
			panic(err.Error())
		}

		InfomationTable = append(InfomationTable, information)
	}

	defer results.Close()
	return InfomationTable
}

func (db dbConn) getNewid() uint64 {

	var Number uint64

	rowNumber, err := db.db.Query("SELECT COUNT(*) FROM  statistic")
	if err != nil {
		panic(err.Error())
	}

	for rowNumber.Next() {

		// for each row, scan the result into our image composite object
		err := rowNumber.Scan(&Number)
		if err != nil {
			panic(err.Error())
		}

	}
	newId := Number + 1
	defer rowNumber.Close()

	return newId
}

func (db dbConn) dataBaseInserion(pod Pod, startScheduling time.Time, timePodMatched time.Time) {

	timeStartRunning := db.getStartRunningTime(*pod.Status.StartTime)
	waitingSchDuration := calculDuration(*pod.Status.StartTime, startScheduling)
	SchDuration := calculDuration(startScheduling, timePodMatched)
	LatencyTime := waitingSchDuration + SchDuration
	creationPodDuration := calculDuration(timePodMatched, timeStartRunning)
	runningPodDuration := calculDuration(timeStartRunning, *pod.Status.Conditions[1].LastTransitionTime)
	deploymentTime := creationPodDuration + runningPodDuration
	totalDuration := calculDuration(*pod.Status.StartTime, *pod.Status.Conditions[1].LastTransitionTime)
	id := db.getNewid()
	insert, err := db.db.Query("INSERT INTO statistic  VALUES ( ?, ?,?,?,?,?,?,?,?,? ,? ,? ,? ,?,?,? )", id, pod.Spec.Containers[0].Image, pod.Spec.NodeName, pod.Status.Phase, pod.Status.StartTime, startScheduling, waitingSchDuration, timePodMatched, SchDuration, LatencyTime, timeStartRunning, creationPodDuration, pod.Status.Conditions[1].LastTransitionTime, runningPodDuration, deploymentTime, totalDuration)

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()

	/*

	   	id:=155

	   	insert, err := db.db.Prepare("INSERT INTO informationtable VALUES ( ?, ?,?,? )")

	   	// if there is an error inserting, handle it
	   	if err != nil {
	   		panic(err.Error())
	   	}
	   	// be careful deferring Queries if you are using transactions


	   	res,err:=insert.Exec(id,ImagePending,node,duration)


	       lastId,err:=res.LastInsertId()

	       if err!=nil {

	       	log.Fatal(err)
	   	}
	   rowCnt,err:=res.RowsAffected()

	   	if err!=nil {

	   		log.Fatal(err)
	   	}



	   fmt.Println(lastId)
	   fmt.Println(rowCnt)

	   	defer insert.Close()
	*/
}

func calculDuration(time1 time.Time, time2 time.Time) float64 {

	t1 := time1
	t2 := time2

	diff := t2.Sub(t1)

	duration := float64(diff.Seconds())

	return duration

}

func addPodtoDB(initial int, final int, Pods *PodList, PendingTime time.Time, timePodMatched time.Time) {

	if final > initial {
		diff := final - initial

		for i := 0; i < diff; i++ {

			db.dataBaseInserion(Pods.Items[i], PendingTime, timePodMatched)

		}

	}

}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//Running pod database

func (db dbConn) dataBaseInserionToRunnninTable(pod Pod, timeStartRunning time.Time) {

	id := db.getNewidRunning()

	durationEstimated := getDurationEstimated(pod)

	FinalTimeEstimated := getFinalTimeEstimated(pod, durationEstimated)
	insert, err := db.db.Query("INSERT INTO runningtableestimation VALUES (?,  ?,?,?,? ,?,?)", id, pod.Spec.Containers[0].Image, pod.Spec.NodeName, pod.Status.StartTime, durationEstimated, FinalTimeEstimated, timeStartRunning)

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()

}

////////////////////////////

func getDurationEstimated(pod Pod) float64 {

	var InformationUsed dataBaseInformationTable
	Information := db.informationSearch(pod.Spec.Containers[0].Image)
	for _, k := range Information {
		if k.NodeUsed == pod.Spec.NodeName && k.ImageUsed == pod.Spec.Containers[0].Image {
			InformationUsed = k
			break
		}
	}

	durationEstimated := InformationUsed.RunningTime

	return durationEstimated

}

//////////////////////////////////////

func getFinalTimeEstimated(pod Pod, durationEstimated float64) time.Time {

	MdurationEstimated := time.Duration(durationEstimated)

	FinalTimeEstimated := (pod.Status.StartTime).Add(MdurationEstimated * time.Second)

	return FinalTimeEstimated
}

/////////////////////
func addPodtoDBrunning(initial int, final int, Pods *PodList, timeStartRunning time.Time) {

	if final > initial {
		diff := final - initial

		for i := 0; i < diff; i++ {

			db.dataBaseInserionToRunnninTable(Pods.Items[i], timeStartRunning)

		}

	}

}

func (db dbConn) getNewidRunning() uint64 {

	var Number uint64

	rowNumber, err := db.db.Query("SELECT COUNT(*) FROM  runningtableestimation ")
	if err != nil {
		panic(err.Error())
	}

	for rowNumber.Next() {

		// for each row, scan the result into our image composite object
		err := rowNumber.Scan(&Number)
		if err != nil {
			panic(err.Error())
		}

	}
	newId := Number + 1
	defer rowNumber.Close()

	return newId
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (db dbConn) getStartRunningTime(StartTime time.Time) time.Time {

	var StartRunning time.Time

	rowStartRunning, err := db.db.Query("SELECT timeStartRunning FROM  runningtableestimation where startTime =?", StartTime)
	if err != nil {
		panic(err.Error())
	}

	for rowStartRunning.Next() {

		err := rowStartRunning.Scan(&StartRunning)
		if err != nil {
			panic(err.Error())
		}

	}

	defer rowStartRunning.Close()

	return StartRunning
}
