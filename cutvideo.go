package cutv

import (
    "os/exec"
    "time"
    "strings"
    "bytes"
    "path/filepath"
    "strconv"
    "fmt"
    "errors"
    "os"
)

func getEpoch( filePath string ) (*time.Time,error){

	fileName := filepath.Base(filePath)
        fileExt := filepath.Ext(filePath)
        didEpochStr := fileName[:len(fileName)-len(fileExt)]
        
	didEpochSlice := strings.Split(didEpochStr,"-")
	if len(didEpochSlice) != 2 {
            err := errors.New( fmt.Sprintf("file:%s is invalid",fileName )  )
            Error(err)
            return nil,err
        }

	epochInt,err := strconv.Atoi( didEpochSlice[1] )
	if err != nil {
	    Error(err)
            return nil,err
	}
	fileEpoch := time.Unix(int64(epochInt),0)

	return &fileEpoch,nil

}

func ffmpegCmd( command string )error{
        Infoln(command)
        commands := strings.Split(command," ")
//	cmd := exec.Command( "/usr/bin/ffmpeg", "-ss", fromStr, "-i", file , "-to" ,toStr, "-c" ,"copy", "cut.mp4" )
	cmd := exec.Command( commands[0],commands[1:]... )
	cmd.Stdin = strings.NewReader("")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
	    Error(err)
            return err
	}
        Info( "ffmpeg:"+out.String() )
        return nil
}

func cut( file string , start_time time.Time, duration int )(string,error){

 	fileEpoch,err := getEpoch(file) 
	if err != nil {
            return "",nil
	}
        Infof("file time:%v", *fileEpoch)
        Infof("start time:%v", start_time )

        offset := start_time.Sub( *fileEpoch )
        Infof("offset:%v", offset )
	fromStr := fmt.Sprintf("00:%02d:%02d", int(offset.Minutes()) , int(offset.Seconds()) )
	toStr := fmt.Sprintf("00:%02d:%02d", duration/60 , duration%60 )

        Infof("from:%s to:%s",fromStr,toStr)
        command := fmt.Sprintf("ffmpeg -ss %s -i %s -to %s -c copy cut.mp4",fromStr,file,toStr )
        err = ffmpegCmd(command)
	if err != nil {
	    Error(err)
            return "",nil
	}

	return "cut.mp4",nil
}


func concatnate( files []string )(string,error){

    f, err := os.OpenFile( "concat.list" , os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
    if err != nil {
        Error(err)
        return "",err
    }

    defer f.Close()
    
    for _,file := range files {
        fmt.Fprintf(f, "file '%s'\n" , file )
    }   


 
    return "concat.list",nil


}

//根據 getRecordFilePath 拿回的 file lists 直接傳到此func的file_path
//根據 start_time,duration 在此組出 ffmpeg cmd 
//轉檔mp4後回傳檔名 若有err則回傳
func GenMP4File(file_lists []string,start_time time.Time,duration int) (string,error){

    Infof("Get %s,cut at:%d , duration:%d",file_lists,start_time,duration)

    var cutFile string
    if len(file_lists)> 1 {
        f,err := concatnate(file_lists)
        cutFile = f
    if err != nil {
    	Error(err)
        return "",err
    }
    }else if len(file_lists)==1{
        cutFile = file_lists[0]
    }
    
    return cut(cutFile,start_time,duration)

}




 
