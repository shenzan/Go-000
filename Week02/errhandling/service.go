package main

import (
	//service层不需要再依赖database/sql
	//"database/sql"
	"fmt"
)

type IErrorNotFound interface {
	error
	IsNotFound() bool
	ErrorMsg() string
}

func getUserName(id int) (string, error) {
	name, err := findUserNameById(id)
	//在service层只识别IErrorNotFound，不需要识别sql.ErrNoRows，和底层实现解耦
	if e, ok := err.(IErrorNotFound); ok && e.IsNotFound() {
		fmt.Printf("error msg from dao: %+v\n", e.ErrorMsg())
		//假设业务需求是找不到用户就返回404, 不算异常
		return "404", nil
	}
	return name, err
}

func main() {
	name, err := getUserName(1)
	if err != nil {
		fmt.Printf("got err. %+v\n", err)
	} else {
		fmt.Printf("got user name. %+v\n", name)
	}

}
