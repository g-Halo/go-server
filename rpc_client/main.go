package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Params struct {
	Username string
	Password string
}

func main() {
	// client, err := rpc.Dial("tcp", ":7072")
	// defer client.Close()
	// if err != nil {
	// 	fmt.Println(err)
	// 	// fmt.Println("无效的地址")
	// 	os.Exit(0)
	// }

	// fmt.Println("调用 token create")
	// p := &Params{Username: "test1", Password: "123"}
	// res := &util.Response{}
	// if err := client.Call("Token.Create", p, &res); err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(res)

	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }

	song := make(map[string]interface{})
	song["username"] = "test2"
	song["password"] = "123"
	bytesData, err := json.Marshal(song)
	reader := bytes.NewReader(bytesData)
	url := "http://localhost:7834/v1/login"
	request, err := http.NewRequest("POST", url, reader)
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
