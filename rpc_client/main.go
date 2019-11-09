package main

import (
	"fmt"
	"github.com/g-Halo/go-server/util"
	"net/rpc"
	"os"
)

type Params struct {
	Username string
	Password string
}

func main() {
	client, err := rpc.Dial("tcp", ":7301")
	if err != nil {
		fmt.Println("无效的地址")
		os.Exit(0)
	}


	//str := "abcd"
	//var reply int
	//err = client.Call("Token.Validate", &str, &reply)
	//fmt.Println(reply)

	fmt.Println("调用 token create")
	p := &Params{Username: "test1", Password: "123"}
	res := &util.Response{}
	client.Call("Token.Create", p, &res)
	fmt.Println(res)

	if err != nil {
		fmt.Println(err.Error())
		return
	}


}
