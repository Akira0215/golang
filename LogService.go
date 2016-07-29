package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

type LogData struct {
	Id     string
	Type   string
	Msg    string
	Option string
	Date   string
}

type LogData2 struct {
	Id      string
	Date    string
	LogText []string
}

// func sayhelloName(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm() //解析url传递的参数，对于POST则解析响应包的主体（request body）
// 	//注意:如果没有调用ParseForm方法，下面无法获取表单的数据
// 	fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
// 	fmt.Println("path", r.URL.Path)
// 	fmt.Println("scheme", r.URL.Scheme)
// 	fmt.Println(r.Form["url_long"])
// 	for k, v := range r.Form {
// 		fmt.Println("key:", k)
// 		fmt.Println("val:", strings.Join(v, ""))
// 	}
// 	fmt.Fprintf(w, "Hello astaxie!") //这个写入到w的是输出到客户端的
// }

// func login(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("method:", r.Method) //获取请求的方法
// 	if r.Method == "GET" {
// 		t, _ := template.ParseFiles("login.gtpl")
// 		t.Execute(w, nil)
// 	} else {
// 		//请求的是登陆数据，那么执行登陆的逻辑判断
// 		fmt.Println("username:", r.Form["username"])
// 		fmt.Println("password:", r.Form["password"])
// 	}
// }

func logRead(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)

	fmt.Println("date: %s\n", r.FormValue("date"))
	fmt.Println("id: %s\n", r.FormValue("id"))

	readFile(w, r, r.FormValue("date"), r.FormValue("id"))

}

// func readFile(w http.ResponseWriter, r *http.Request, dirname string, filename string) {

// 	var filePath = "logs" + "/" + dirname + "/" + filename + ".txt"
// 	fmt.Println("filePath: %s\n", filePath)
// 	if Exist(filePath) {
// 		fmt.Println("有此檔案")

// 		f, _ := os.Open(filePath)
// 		// Create a new Scanner for the file.
// 		scanner := bufio.NewScanner(f)

// 		fmt.Println(scanner.Text())

// 		fmt.Fprintf(w, scanner.Text())

// 	} else {
// 		fmt.Fprintf(w, "沒有相關log檔案")
// 	}

// }

func readFile(w http.ResponseWriter, r *http.Request, dirname string, filename string) {

	var filePath = "logs" + "/" + dirname + "/" + filename + ".txt"

	if Exist(filePath) {

		var lines []string
		// var err error

		lines, _ = readLines(filePath)

		var mData = LogData2{r.FormValue("id"), r.FormValue("date"), lines}

		var b, _ = json.Marshal(mData)

		fmt.Fprintf(w, string(b))

	} else {
		fmt.Fprintf(w, "沒有相關log檔案")
	}

}

func readLines(path string) (lines []string, err error) {
	var (
		file   *os.File
		part   []byte
		prefix bool
	)
	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0))
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			lines = append(lines, buffer.String())
			buffer.Reset()
		}
	}

	// if err == io.EOF {
	// 	err = nil
	// }
	return
}

func logWrite(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)

	fmt.Println(r.Form["id"])
	fmt.Println(r.Form["type"])
	fmt.Println(r.Form["msg"])
	fmt.Println(r.Form["option"])
	// for k, v := range r.Form {
	//     fmt.Println("key:", k)
	//     fmt.Println("val:", strings.Join(v, ""))
	// }

	t := time.Now()
	fmt.Println(t) // e.g. Wed Dec 21 09:52:14 +0100 RST 2011
	// fmt.Printf("%02d.%02d.%4d\n", t.Day(), t.Month(), t.Year())
	var sDate = fmt.Sprintf("%4d%02d%02d", t.Year(), t.Month(), t.Day())
	var sTime = fmt.Sprintf("%4d%02d%02d%02d%02d%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

	var mData = LogData{r.FormValue("id"), r.FormValue("type"), r.FormValue("msg"), r.FormValue("option"), sTime}

	var b, _ = json.Marshal(mData)

	var filename = r.FormValue("id") + ".txt"

	go writeFile(sDate, filename, string(b))

	fmt.Fprintf(w, string(b)) //这个写入到w的是输出到客户端的
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)

	fmt.Println("runtime.NumCPU() : ", runtime.NumCPU())

	// http.HandleFunc("/", sayhelloName)
	http.HandleFunc("/log/write", logWrite)
	http.HandleFunc("/log/read", logRead)
	http.HandleFunc("/log/test", test)
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func writeFile(dirname string, filename string, data string) {

	checkFile(dirname, filename)

	var filePath = "logs" + "/" + dirname + "/" + filename

	if Exist(filePath) {

		outputf, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0600)
		//                 ^^^^^^^^                ^^^^^^^^^^^^^^^^^^      ^^^^
		// 是用 OpenFile,不是只用Open,因為還要設定模式. 建立檔案 只有寫入  UNIX檔案權限
		if err != nil {
			fmt.Println("開檔錯誤!")
			return
		}
		defer outputf.Close() // 離開時關檔

		// outStr := "Golang寫檔測試\n"
		outputWriter := bufio.NewWriter(outputf) // 建立緩衝輸出物件

		outputWriter.WriteString(data + "\r\n")

		outputWriter.Flush()

	}

}

func checkFile(dirname string, filename string) {

	fmt.Println(dirname)
	fmt.Println(filename)

	if !Exist("logs") {
		os.Mkdir("logs", os.ModePerm)
	}

	var path = "logs/" + dirname

	fmt.Println("path : ", path)

	if !Exist(path) {
		fmt.Println("create path", path)
		os.MkdirAll(path, os.ModePerm)
	} else {
		fmt.Println("path is exist", path)
	}

	var filepath = path + "/" + filename

	if !Exist(filepath) {
		fmt.Println("create file : ", filepath)
		os.Create(filepath)
	} else {
		fmt.Println("file is exist : ", filepath)
	}

}

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func test(w http.ResponseWriter, r *http.Request) {
	// go func() {
	// t := time.Now()
	fmt.Fprintf(w, "1234567890")
	// }()

}
