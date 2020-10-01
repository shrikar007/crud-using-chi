package common

import (
	"encoding/csv"
	"os"
	"time"
)

func GenerateCSV(fileName string, data [][]string) error{
	var (
		path string
		writer *csv.Writer
		file *os.File
		err error
	)

	path = "export/"+GetFileName(fileName)+".csv"
	if file, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.ModePerm); err!=nil{
		return err
	}
	defer file.Close()

	defer func() {
		if writer != nil {
			writer.Flush()
		}
	}()
	writer = csv.NewWriter(file)
	if err = writer.WriteAll(data); err!=nil{
		return err
	}
	return err
}

func GetFileName(fileName string) string{
	var timeStamp = time.Now().Format(time.RFC3339)
	return fileName+"_"+timeStamp
}
