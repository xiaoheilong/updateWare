package softwareVersionManager

import (
	"archive/zip"
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func Unzip(zipFile string, destDir string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		log.Println("openerr:" + err.Error())
		return err
	}
	defer zipReader.Close()
	for _, f := range zipReader.File {
		fpath := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				log.Println("mkdirerr:" + err.Error())
				return err
			}
			inFile, err := f.Open()
			if err != nil {
				log.Println("infileopen:" + err.Error())
				return err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				log.Println("outfileopen:" + err.Error())
				return err
			}
			defer outFile.Close()
			_, err = io.Copy(outFile, inFile)
			if err != nil {
				log.Println("copyerr:" + err.Error())
				return err
			}
		}
	}
	return nil
}

func init() {
	logFile, err := os.OpenFile("log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	log.SetOutput(logFile)
	log.SetPrefix("[update]")
}

func main() {
	var inFile, outFile, exeName string

	flag.StringVar(&inFile, "zip-path", "", "zip-path")
	flag.StringVar(&outFile, "dest-path", "", "dest-path")
	flag.StringVar(&exeName, "app-name", "", "app-name")

	flag.Parse()

	if inFile == "" {
		return
	}

	log.Println("inFile:" + inFile)
	log.Println("outFile:" + outFile)
	log.Println("exe" + exeName)
	//防止过快资源没释放
	time.Sleep(time.Duration(1) * time.Second)
	unzipErr := Unzip(inFile, outFile)
	if unzipErr == nil {

		datapath := filepath.Join(outFile, exeName)
		log.Println("datapath:" + filepath.ToSlash(datapath))
		//cmd := exec.Command("cmd.exe", "/c", command)
		cmd := exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", filepath.ToSlash(datapath))
		cmd.Run()
	}
}
