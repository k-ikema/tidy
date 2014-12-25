package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"code.google.com/p/gcfg"
)

const RFC1123 = "2006/01/02 15:04:05 MST"
const IniFile = "tidy.ini"

type config struct {
	Path struct {
		Workpath string
		Rulepath string
	}
}

func main() {

	IniFilePath, WorkPath, RulePath, err := initPath()

	if "" != err {
		fmt.Println(err)
		pause()
		os.Exit(0)
	}

	fmt.Println(IniFilePath)
	fmt.Println(WorkPath)
	fmt.Println(RulePath)

	pause()
}

func getRuleFilePath() (string, bool) {
	var RuleFilePath string
	if len(os.Args) == 1 {
		ExecPath, _ := os.Getwd()
		RuleFilePath = ExecPath + "\\" + createDefaultRuleFileName()
	} else {
		RuleFilePath = os.Args[1]
	}
	_, err := os.Stat(RuleFilePath)
	return RuleFilePath, !os.IsNotExist(err)
}

func createDefaultRuleFileName() string {
	t := strings.Split(time.Now().Format(RFC1123), "/")
	year := t[0]
	quarter, _ := strconv.Atoi(t[1])
	quarter = (quarter-1)/3 + 1
	return (year + "Q" + strconv.Itoa(quarter) + ".txt")
}

func checkArgs() bool {
	return true
}

func showUsage() {
	fmt.Println("Usage : file_dist.exe [ini_File's_Path]")
}

func pause() {
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func initPath() (string, string, string, string) {
	//var IniFilePath string
	var PathError string
	var ErrorMessage string

	IniFilePath, err := getIniFile()
	if false == err {
		ErrorMessage += "ini file is not exist. \n"
		return IniFilePath, "", "", ErrorMessage
	}
	//iniファイル読み取り
	WorkPath, RulePath, PathError := readIniFile(IniFilePath)

	if "" != PathError {
		ErrorMessage += PathError
		return IniFilePath, "", "", ErrorMessage
	}

	return IniFilePath, WorkPath, RulePath, ErrorMessage
}

// ini ファイルのpath 取得と存在チェック
// 引数でini ファイルのパスが指定されていない場合は、実行ディレクトリ内をチェック
func getIniFile() (string, bool) {
	var IniFilePath string
	if len(os.Args) == 1 {
		ExecPath, _ := os.Getwd()
		IniFilePath = ExecPath + "\\" + IniFile
	} else {
		IniFilePath = os.Args[1]
	}
	_, err := os.Stat(IniFilePath)
	return IniFilePath, !os.IsNotExist(err)
}

// ini ファイルの内容読み取り
func readIniFile(IniFilePath string) (string, string, string) {

	var cfg config
	var ErrorMessage string
	var RulePath string

	err := gcfg.ReadFileInto(&cfg, IniFilePath)

	if err != nil {
		log.Fatalf("Failed to parse config file: %s", err)
	}
	//必須項目 workpath のチェック
	if "" == cfg.Path.Workpath {
		ErrorMessage += "No setting for Path to Work Folder.\n"
	}
	//ワークフォルダ存在チェックとディレクトリ判定
	WorkPath := cfg.Path.Workpath
	WorkPathInfo, err := os.Stat(WorkPath)
	if false == !os.IsNotExist(err) {
		ErrorMessage += "Path '" + WorkPath + "' is not Exist.\n"
	} else if false == WorkPathInfo.IsDir() {
		ErrorMessage += "Path '" + WorkPath + "' is not Directory.\n"
	}

	//Rulepath 設定値のフォルダ判定
	//　フォルダあるいは設定無しだった場合のデフォルトファイル名でのパス設定
	RulePath = cfg.Path.Rulepath
	if "" != RulePath {
		PathInfo, _ := os.Stat(RulePath)
		if true == PathInfo.IsDir() {
			RulePath += "\\" + createDefaultRuleFileName()
		}
	} else {
		ExecPath, _ := os.Getwd()
		RulePath += ExecPath + "\\" + createDefaultRuleFileName()
	}

	//　規則ファイルの存在チェック
	_, err = os.Stat(RulePath)
	if false == !os.IsNotExist(err) {
		ErrorMessage += "Rule File is not exist.　-> " + RulePath + "\n"
	}

	return WorkPath, RulePath, ErrorMessage
}
