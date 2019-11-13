package main

import (
	"fmt"
	"net/rpc"
	"os"

	"github.com/g-Halo/go-server/util"
)

type Params struct {
	Username string
	Password string
}

func main() {
	client, err := rpc.Dial("tcp", ":7072")
	defer client.Close()
	if err != nil {
		fmt.Println(err)
		// fmt.Println("无效的地址")
		os.Exit(0)
	}

	fmt.Println("调用 token create")
	p := &Params{Username: "test1", Password: "123"}
	res := &util.Response{}
	if err := client.Call("Token.Create", p, &res); err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

}
