package main

import (
	"time"
	"cutv"
)

func main(){

    t := time.Unix( 1483620742,0 )
    //file,err := genMP4File( []string{"/mnt/rec/923bb4ec937c07324de968025483f617-1483606525.flv"},t ,30)
    file,err := cutv.GenMP4File( []string{
	    "/mnt/rec/923bb4ec937c07324de968025483f617-1483620732.flv",
	    "/mnt/rec/923bb4ec937c07324de968025483f617-1483620772.flv",
	    "/mnt/rec/923bb4ec937c07324de968025483f617-1483620813.flv",
    },t ,150)

    if err != nil {
    	cutv.Error(err)
        return
    }

    cutv.Infof("cut file at %s",file)

}


 
