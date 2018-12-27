
package main

import (
	"database/sql"
	"k8s.io/client-go/kubernetes"
	_ "kubCG/go-sql-driver/mysql"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const schedulerName = "hightower"


type dbConn struct {
	db *sql.DB
}

type ConfigFile struct {

	kubeconfig *string
	clientset *kubernetes.Clientset
	err *error


}


var db dbConn


var  config ConfigFile

func main() {

	log.Println("Starting custom scheduler...")


	config.openFile()

	var err error
	db.db, err = sql.Open("mysql", "gofrane:gofrane@tcp(127.0.0.1:3306)/kubcgdatabase")

	if err != nil {
		log.Println("The scheduler can not connect to the database server !! ")
	}





	defer db.db.Close()

	doneChan := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(1)
	go monitorUnscheduledPods(doneChan, &wg)

	wg.Add(1)
	go reconcileUnscheduledPods(30, doneChan, &wg)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-signalChan:
			log.Printf("Shutdown signal received, exiting...")
			close(doneChan)
			wg.Wait()
			os.Exit(0)
		}
	}
}