package main

import (
lazlo "github.com/djosephsen/lazlo/lib"
"os/signal"
"syscall"
)

func main(){
	
	//make a broker
	broker,err := lazlo.NewBroker()
	if err != nil{
		lazlo.Logger.Error(err)
		return
	}
	brain := *broker.Brain
	defer brain.Close()

	//start the read, write and broker threads
	broker.Start()

	if err := broker.InitModules(); err !=nil{
      lazlo.Logger.Error(err)
		return
	}

	// Loop
	signal.Notify(broker.SigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
  	stop := false
   for !stop {
      select {
      case sig := <-broker.SigChan:
         switch sig {
         case syscall.SIGINT, syscall.SIGTERM:
            stop = true
         }
      }
   }
   // Stop listening for new signals
   signal.Stop(broker.SigChan)

	//wait for the write thread to stop (so the shutdown hooks have a chance to run)
	broker.WriteThread.RunChan <- true
	<- broker.SyncChan
}