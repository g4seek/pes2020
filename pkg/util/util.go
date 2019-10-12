package util

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

const BaseDir string = "./src/pes2020/data/"

func RenewFile(fileName string) {
	err := os.RemoveAll(BaseDir + "/" + fileName)
	if err != nil {
		panic(err)
	}
	file, err := os.Create(BaseDir + "/" + fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
}

func AppendLine(fileName string, line string) {
	f, err := os.OpenFile(BaseDir+"/"+fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		fmt.Println("os OpenFile error: ", err)
		return
	}
	defer f.Close()
	_, err = f.WriteString(line)
	if err != nil {
		panic(err)
	}
}

func ReadLines(fileName string) []string {
	filePath, _ := filepath.Abs(BaseDir + "/" + fileName)
	fi, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return nil
	}
	defer fi.Close()
	// 逐行读取文件
	br := bufio.NewReader(fi)
	var lines []string
	for {
		lineInByte, _, e := br.ReadLine()
		if e == io.EOF {
			break
		}
		line := string(lineInByte)
		lines = append(lines, line)
	}
	return lines
}

func GetRequest(url string, headers, cookies map[string]string) (response string) {
	client := http.Client{}

	// 创建get请求
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	// 添加cookie信息
	for key, value := range cookies {
		request.AddCookie(&http.Cookie{Name: key, Value: value})
	}
	// 添加header信息
	for key, value := range headers {
		request.Header.Add(key, value)
	}

	resp, err := client.Do(request)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}
	response = result.String()
	return
}

func ParseStrToInt(str string) int {
	value, err := strconv.Atoi(str)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	return value
}
