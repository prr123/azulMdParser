package main


import (
	"fmt"
	"log"
	"os"
	"bytes"

	"github.com/goccy/go-yaml"
	util "github.com/prr123/utility/utilLib"

)

func main() {

    numarg := len(os.Args)
    flags:=[]string{"dbg", "azmd"}

    useStr := "/azmd=azmdfile [/dbg]"
    helpStr := "program that parses azmd files"


    if numarg == 1 {
        fmt.Printf("help: %s\n", helpStr)
        fmt.Printf("usage is: %s %s\n", os.Args[0], useStr)
        os.Exit(1)
    }

    if numarg > len(flags) + 1 {
        fmt.Println("too many arguments in cl!")
        fmt.Println("usage: %s %s\n", os.Args[0], useStr)
        os.Exit(1)
    }

    if numarg > 1 && os.Args[1] == "help" {
            fmt.Printf("help: %s\n", helpStr)
            fmt.Printf("usage is: %s %s\n", os.Args[0], useStr)
            os.Exit(1)
	}

    flagMap, err := util.ParseFlags(os.Args, flags)
    if err != nil {log.Fatalf("error -- parsing cli flags: %v\n", err)}

    dbg:= false
    _, ok := flagMap["dbg"]
    if ok {dbg = true}

    azmdStr := ""
   	azmdval, ok := flagMap["azmd"]
    if ok {
        if azmdval.(string) == "none" {log.Fatalf("error -- no azmd file provided!\n")}
        azmdStr = azmdval.(string)
    }

	err = parseAzmd(azmdStr, dbg)
	if err != nil {log.Fatalf("error -- parsing azmd: %v\n", err)}

	log.Printf("Success!\n")
}

type head struct {
	Title string `yaml:"title"`
	Author string `yaml:"author"`
	Short string `yaml:"short"`
	DateStr string `yaml:"date"`
	Keys []string `yaml:"keys"`
}

func parseAzmd(azmdStr string, dbg bool) (err error) {

	azmdFilnam := "azulMd/" + azmdStr + ".azmd"
	if dbg {
		fmt.Printf("dbg -- azmd File: %s\n", azmdFilnam)
	}

	azmd, err := os.ReadFile(azmdFilnam)
	if err != nil {return fmt.Errorf("error -- read azmd: %v\n", err)}

	if dbg {
		fmt.Printf("***** %s ******\n", azmdFilnam)
		fmt.Printf("%s\n", azmd)
		fmt.Println("******** end ***********")
	}

	sep := []byte("---\n")
	yamlSt := bytes.Index(azmd,sep)
	if yamlSt == -1 {return fmt.Errorf("no yaml header found!")}

	yamlEnd := bytes.Index(azmd[yamlSt+4:], sep)
	if yamlEnd == -1 {return fmt.Errorf("no yaml end found!")}
	yamlEnd += yamlSt+4
	yamlSect := azmd[yamlSt+4:yamlEnd]
//	fmt.Printf("%d:%d\n%s", yamlSt, yamlEnd, azmd[yamlSt:yamlEnd])
	if dbg {
		fmt.Println("***** yaml section ******\n")
		fmt.Printf("%s\n", yamlSect)
		fmt.Println("******** end ***********")
	}

	sumTit := []byte("#Summary\n")
	sumSt := bytes.Index(azmd[yamlEnd+4:],sumTit)
	if yamlSt == -1 {return fmt.Errorf("no summary begin found!")}
	sumSt += yamlEnd + 4 +9
	sumEnd := bytes.Index(azmd[sumSt:], []byte("\n#"))
	if sumEnd == -1 {return fmt.Errorf("no summary end found!")}
	sumEnd += sumSt
	sumSect := azmd[sumSt:sumEnd]

	//	sumSect
	if dbg {
		fmt.Println("***** summary section ******\n")
		fmt.Printf("%s\n", sumSect)
		fmt.Println("******** end ***********")
	}

	// main Section
	mainSt := sumEnd +1
	mainSect := azmd[mainSt:]
	if dbg {
		fmt.Println("***** main section ******\n")
		fmt.Printf("%s\n", mainSect)
		fmt.Println("******** end ***********")
	}

	// validate yaml
	var inp head
	err = yaml.Unmarshal(yamlSect, &inp)
	if dbg {
		fmt.Printf("****** Head ******\n")
		fmt.Printf("  Title: %s\n",inp.Title)
		fmt.Printf("  Author: %s\n",inp.Author)
		fmt.Printf("  Short: %s\n",inp.Short)
		fmt.Printf("  Date: %s\n",inp.DateStr)
		fmt.Printf("  Keys:\n")
		for i:=0; i< len(inp.Keys); i++ {
			fmt.Printf("   %d: %s\n", i+1, inp.Keys[i])
		}
		fmt.Printf("**** End Head ****\n")
	}

	//create dir if it does not exist
	_, err = os.Stat("articles/" + azmdStr)
	if os.IsNotExist(err) {
		log.Printf("info -- found no article dir! creating...!\n")
		err1 := os.Mkdir("articles/" + azmdStr, 0750)
		if err1 != nil {return fmt.Errorf("error -- creating dir: %v!", err1)}
	} else {
		log.Printf("info -- found dir! removing existing files")
		err1 := os.RemoveAll("articles/" + azmdStr)
		if err1 != nil {return fmt.Errorf("error -- removing files: %v!", err1)}
		err1 = os.Mkdir("articles/" + azmdStr, 0750)
		if err1 != nil {return fmt.Errorf("error -- creating dir: %v!", err1)}
	}

	err = os.WriteFile("articles/" + azmdStr + "/head.yaml" , yamlSect, 0666)
	if err != nil {return fmt.Errorf("error -- creating head file: %v!", err)}
	err = os.WriteFile("articles/" + azmdStr + "/summary.md" , sumSect, 0666)
	if err != nil {return fmt.Errorf("error -- creating summary file: %v!", err)}
	err = os.WriteFile("articles/" + azmdStr + "/main.md" , mainSect, 0666)
	if err != nil {return fmt.Errorf("error -- creating main file: %v!", err)}
	return nil
}
