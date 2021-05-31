package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	go forever()

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
}

func forever() {

	for {
		os.Setenv("GOOS", "windows")
		os.Setenv("GOARCH", "amd64")
		fmt.Println("Type in name of series:")
		in := bufio.NewReader(os.Stdin)
		name, err := in.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading string: %v\n", err)
			return
		}
		re := regexp.MustCompile(`\r?\n`)
		name = re.ReplaceAllString(name, "")

		fmt.Println("Type in how many volumes:")
		in = bufio.NewReader(os.Stdin)
		volumes, err := in.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading string: %v\n", err)
			return
		}
		volumes = re.ReplaceAllString(volumes, "")
		volumesInt, err := strconv.ParseInt(volumes, 10, 64)
		if err != nil {
			fmt.Printf("Error converting int %v\n", err)
			return
		}

		createDir(name)
		for i := 1; i <= int(volumesInt); i++ {
			base := fmt.Sprintf("%s/%s_v%d", name, name, i)
			spread := fmt.Sprintf("%s_v%d_spreads", name, i)
			scans := fmt.Sprintf("%s_v%d_scans", name, i)
			createDir(base)
			createDir(base + "/" + spread)
			createDir(base + "/" + scans)
		}

		err = renameScans(name)

		if err != nil {
			fmt.Println("Err processing scans: ", err)
		}

		fmt.Println("\n\n\n\nAll done, to close hit ctrl+c")
		fmt.Println("\nIf you want to continue just redo the previous steps")
		fmt.Println("==================================")
		fmt.Println("==================================")
		fmt.Println("==================================")
	}

}

func createDir(dir string) {
	_, err := os.Stat(dir)

	if os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)

		if err != nil {
			log.Fatal(err)
		}
	}

}

func renameScans(dir string) error {

	dir = fmt.Sprintf("%s/scans", dir)
	folders, err := ioutil.ReadDir(dir)

	if err != nil {
		return err
	}

	err = updateFiles(dir, folders)
	if err != nil {
		return err
	}

	return nil
}

func updateFiles(dir string, folders []os.FileInfo) error {
	for _, fo := range folders {
		files, err := ioutil.ReadDir(fmt.Sprintf("%s/%s", dir, fo.Name()))
		SortFileNameAscend(files)
		if err != nil {
			return err
		}
		for i, fi := range files {
			fileArray := strings.Split(fi.Name(), ".")
			ext := fileArray[len(fileArray)-1]
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			orginalPath := fmt.Sprintf("%s/%s/%s/%s", cwd, dir, fo.Name(), fi.Name())
			newPath := fmt.Sprintf("%s/%s/%s/%s", cwd, dir, fo.Name(), fmt.Sprintf("%s__%d.%s", fo.Name(), i+1, ext))

			os.Rename(orginalPath, newPath)
		}

	}

	return nil
}

func SortFileNameAscend(files []os.FileInfo) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
}
