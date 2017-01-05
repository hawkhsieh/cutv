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
	"path"
)

type FileAttr struct {
    did string
    fileTime time.Time
}


func getFileAttr( filePath string ) (*FileAttr,error) {

	fa := new(FileAttr)
	fileName := filepath.Base(filePath)
	fileExt := filepath.Ext(filePath)
	didEpochStr := fileName[:len(fileName) - len(fileExt)]

	didEpochSlice := strings.Split(didEpochStr, "-")
	if len(didEpochSlice) != 2 {
		err := errors.New(fmt.Sprintf("file:%s is invalid", fileName))
		Error(err)
		return nil, err
	}

	fa.did = didEpochSlice[0]

	epochInt, err := strconv.Atoi(didEpochSlice[1])
	if err != nil {
		Error(err)
		return nil, err
	}

	fa.fileTime = time.Unix(int64(epochInt), 0)
	return fa, nil

}

func ffmpegCmd( command string )error {
	Infoln(command)
	commands := strings.Split(command, " ")
	//	cmd := exec.Command( "/usr/bin/ffmpeg", "-ss", fromStr, "-i", file , "-to" ,toStr, "-c" ,"copy", "cut.mp4" )
	cmd := exec.Command(commands[0], commands[1:]...)
	cmd.Stdin = strings.NewReader("")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		Error(err)
		return err
	}
	Info("ffmpeg:" + out.String())
	return nil
}

func cut( file string , start_time time.Time, duration int )(string,error) {

	fa, err := getFileAttr(file)
	if err != nil {
		return "", nil
	}
	Infof("file time:%v", fa.fileTime)
	Infof("start time:%v", start_time)

	offset := start_time.Sub(fa.fileTime)
	Infof("offset:%v", offset)
	fromStr := fmt.Sprintf("00:%02d:%02d", int(offset.Minutes()), int(offset.Seconds()))
	toStr := fmt.Sprintf("00:%02d:%02d", duration / 60, duration % 60)

	Infof("from:%s to:%s", fromStr, toStr)
	distFile := fmt.Sprintf("%s-%d-%d.mp4",fa.did,fa.fileTime.Unix(),start_time.Unix())
	command := fmt.Sprintf("ffmpeg -y -ss %s -i %s -to %s -c copy %s", fromStr, file, toStr, distFile )
	err = ffmpegCmd(command)
	if err != nil {
		Error(err)
		return "", nil
	}

	os.Remove(file)
	return distFile, nil
}

func concatnate( files []string )(string,error) {

	concatFile := path.Base(files[0])

	concatListFile := strings.Split(concatFile,".")[0] + ".list"

	f, err := os.OpenFile(concatListFile, os.O_WRONLY | os.O_CREATE, 0600)
	if err != nil {
		Error(err)
		return "", err
	}

	defer f.Close()

	for _, file := range files {
		fmt.Fprintf(f, "file '%s'\n", file)
	}

	command := fmt.Sprintf("ffmpeg -y -f concat -safe 0 -i %s -c copy %s",concatListFile , concatFile )
	err = ffmpegCmd(command)
	if err != nil {
		Error(err)
		return "", nil
	}
	os.Remove(concatListFile)

	return concatFile, nil

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




 
