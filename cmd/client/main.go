package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var endpoint string = "http://localhost:8080/"

type ShortenRequest struct {
	URL string `json:"url"`
}

func main() {

	// приглашение в консоли
	fmt.Println("Введите длинный URL")
	// открываем потоковое чтение из консоли
	reader := bufio.NewReader(os.Stdin)
	// читаем строку из консоли
	long, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	long = strings.TrimSuffix(long, "\r\n")
	// заполняем контейнер данными
	// добавляем HTTP-клиент
	client := &http.Client{}
	// пишем запрос
	// запрос методом POST должен, помимо заголовков, содержать тело
	// тело должно быть источником потокового чтения io.Reader
	var request *http.Request

	fmt.Println("Сжать запрос (принять ответ) gzip? y/n")
	// читаем строку из консоли
	gzip, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	gzip = strings.TrimSuffix(gzip, "\r\n")

	isCompress := gzip == "y"

	fmt.Println("Какой запрос отправить? 1-text/plain, 2-json")
	// читаем строку из консоли
	rqType, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	rqType = strings.TrimSuffix(rqType, "\r\n")

	switch rqType {
	case "1":
		request, err = getRequestText(long, isCompress)
	case "2":
		request, err = getRequestJSON(long, isCompress)
	default:
		panic("Select 1 or 2!")
	}

	if err != nil {
		panic(err)
	}

	// отправляем запрос и получаем ответ
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	// выводим код ответа
	fmt.Println("Статус-код ", response.Status)
	defer response.Body.Close()
	// читаем поток из тела ответа
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	// и печатаем его
	fmt.Println(string(body))
}

func getRequestText(URL string, isCompress bool) (*http.Request, error) {
	request, err := getBaseRequest(URL, isCompress, "")
	if err != nil {
		return nil, err
	}
	// в заголовках запроса указываем кодировку
	request.Header.Add("Content-Type", "text/plain")

	return request, nil
}

func getRequestJSON(URL string, isCompress bool) (*http.Request, error) {
	rqModel := ShortenRequest{URL: URL}
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(rqModel)

	request, err := getBaseRequest(buf.String(), isCompress, "api/shorten")
	if err != nil {
		return nil, err
	}
	// в заголовках запроса указываем кодировку
	request.Header.Add("Content-Type", "application/json")

	return request, nil
}

func getBaseRequest(data string, isCompress bool, path string) (*http.Request, error) {
	reader, err := getReader(data, isCompress)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, endpoint+path, *reader)
	if err != nil {
		return nil, err
	}
	if isCompress {
		request.Header.Add("Content-Encoding", "gzip")
		//request.Header.Add("Accept-Encoding", "gzip")
	}

	return request, nil
}

func getReader(data string, isCompress bool) (*io.Reader, error) {
	var rqReader io.Reader

	if isCompress {
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(data))
		defer zb.Close()
		if err != nil {
			return nil, err
		}
		rqReader = buf
	} else {
		rqReader = strings.NewReader(data)
	}

	return &rqReader, nil
}
