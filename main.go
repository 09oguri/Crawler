// Crawler project main.go
package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {
	// WEBページを取得する間隔
	fetchInterval := 3 * time.Second

	// 引数はURLが行区切りで書かれたテキストファイルのパス
	if len(os.Args) < 2 {
		fmt.Println("usage: Crawler [textfile]")
		return
	}

	fp, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	reader := bufio.NewReader(fp)

	// 1ページずつ取得して保存
	for {
		url, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		urlStr := string(url[:])

		// 空行など URL の先頭が http でない場合は無視する
		if !strings.HasPrefix(urlStr, "http") {
			continue
		}

		isCreated := CreateOutputFile(urlStr)

		// WEBページを取得できなかった場合は終了する
		if isCreated < 0 {
			break
		}

		time.Sleep(fetchInterval)
	}

}

// WEBページを取得して保存する
func CreateOutputFile(url string) int {
	output := ExtractFileName(url)

	body := FetchBody(url)
	if body == nil {
		fmt.Println("Cannot fetch the URL: ", url)
		return -1
	}

	ioutil.WriteFile(output, body, os.ModePerm)
	fmt.Println("output: ", output)
	return 0
}

// URLからファイル名を抽出する
func ExtractFileName(url string) string {
	// TODO ファイル名がない場合、拡張子がない場合
	expr := "https?://([a-z|A-Z|0-9|.%-_]+/)*"

	reg, err := regexp.Compile(expr)
	if err != nil {
		panic(err)
	}
	filename := reg.ReplaceAllString(url, "")
	filename = ReplaceInvalidChars(filename)

	return filename
}

// WEBページを取得する
func FetchBody(url string) []byte {
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	return body
}

// ファイル名に使用禁止文字が含まれている場合は _ に置換する
func ReplaceInvalidChars(filename string) string {
	expr := "[/|>|<|?|:|\"|\\|/*|/||;]"

	reg, err := regexp.Compile(expr)
	if err != nil {
		panic(err)
	}
	filename = reg.ReplaceAllString(filename, "_")

	return filename
}
