package main

import (
	"github.com/XinRoom/dir2md/golimit"
	"flag"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"reflect"
	"strings"
	"sync"
)

var (
	PacDocPath = ""
	AllowExt = []string{".docx", ".html", ".doc"}
	OutPath = ""
	NeedConFiles = []string{}
)

//type Charset string
//const (
//	UTF8    = Charset("UTF-8")
//	GB18030 = Charset("GB18030")
//)

// 判断元素是否在对应类型的数组中（通用型）
func IsContain(item interface{}, items interface{}) bool {

	value1 := reflect.ValueOf(item)
	if value1.Kind() == reflect.Array {
		panic("item can not Array!")
		// return false
	}
	newItem := reflect.ValueOf(item)
	// fmt.Printf("item is %s[%s] type \n", newItem.Kind(), newItem.String())

	value2 := reflect.ValueOf(items)
	if value2.Kind() != reflect.Slice {
		panic("items is must Array!")
		// return false
	}
	newItems := make([]reflect.Value, 0, value2.Len())
	for i := 0; i < value2.Len(); i++ {
		_newItem := value2.Index(i)
		// fmt.Printf("items.index %d is %s[%s] type \n", i, _newItem.Kind(), _newItem.String())
		newItems = append(newItems, _newItem)
	}

	for _, eachItem := range newItems {
		if eachItem.Type() == newItem.Type() && eachItem.String() == newItem.String(){
			return true
		}
	}
	return false
}

// 遍历所有的文件
func GetAllFile(pathname string) error {
	rd, err := ioutil.ReadDir(pathname)
	for _, fi := range rd {
		_pathname := path.Join(pathname,fi.Name())
		if fi.IsDir() {
			// fmt.Printf("[%s]\n", path.Join(pathname,fi.Name()))
			err = GetAllFile(_pathname)
		} else if ext := path.Ext(_pathname); IsContain(ext, AllowExt) {
			NeedConFiles = append(NeedConFiles, _pathname)
		}
	}
	return err
}


//func ConvertByte2String(byte []byte, charset Charset) string {
//	var str string
//	switch charset {
//	case GB18030:
//		var decodeBytes,_=simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
//		str= string(decodeBytes)
//	case UTF8:
//		fallthrough
//	default:
//		str = string(byte)
//	}
//	return str
//}

// 调用Pandoc
func ToMd(inFile string, outFile string, args string) error {
	arg := []string{inFile, "--to=markdown_strict+hard_line_breaks+inline_notes", "-o", outFile}
	if path.Ext(inFile) == ".docx" {
		arg = append(arg, "--extract-media=.")
	}
	cmd := exec.Command(PacDocPath, arg...)
	cmd.Dir = path.Dir(outFile)  //设置工作路径，保证相对路径的资源
	// fmt.Println(PacDocPath, arg, cmd.Env)

	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(err.Error(), "file does not exist"){
			panic(err)
		}
		fmt.Printf("Error: [%s] %s\n", inFile, output)
		panic(err)
		return err
	}
	return nil
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

// 协程
func Runc(pathname string) int {

	// 填充待转换文件列表
	_ = GetAllFile(pathname)

	if OutPath == "inPath_out" {
		OutPath = path.Join(pathname, "../", path.Base(pathname) + "_out")
	}

	g := golimit.NewG(10)
	wg := &sync.WaitGroup{}

	numberTasks := len(NeedConFiles)
	if numberTasks == 0 {
		fmt.Printf("Error: numberTasks is 0\n")
		return 1
	}
	bar := progressbar.Default(int64(numberTasks))

	for i := 0; i < numberTasks; i++ {
		wg.Add(1)
		task := NeedConFiles[i]
		g.Run(func() {
			outPath := strings.Replace(strings.TrimSuffix(task, path.Ext(task)) + ".md", pathname, OutPath,1)
			if !FileExist(outPath) {
				_ = os.MkdirAll(path.Dir(outPath), 0644)
				_ = ToMd(task, outPath, "")
			}
			_ = bar.Add(1)
			wg.Done()
		})
	}
	wg.Wait()
	return 0
}

func main() {
	PacDocPath = "pandoc" // 定义pandoc的路径命令或者路径
	var (
		pathname = ""
	)
	flag.StringVar(&pathname, "i", "", "输入文件路径")
	flag.StringVar(&OutPath, "o", "%inPath%_out", "输出文件路径")
	flag.Usage = func () {
		fmt.Println("Usage: dir2md inPath [outPath]")
		flag.PrintDefaults()	// 默认提示信息
	}
	flag.Parse()
	if pathname == "" {
		flag.Usage()
		os.Exit(1)
	}
	if !FileExist(pathname)  {
		fmt.Println(pathname + "is not exist")
		os.Exit(1)
	}
	os.Exit(Runc(pathname))
}