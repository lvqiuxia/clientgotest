package main

import (
	"fmt"
	"practiceCode/clientgotest/clientgo"
	"practiceCode/clientgotest/rest"
)

func main(){

	//对k8s集群进行认证，初始化clientset
	//使用集群外的形式
	mode := clientgo.OutOfClusterClientConfig
	clientSet, err := clientgo.CreateClientSet(mode)
	if err != nil{
		fmt.Println("CreateClientSet fail",err)
	}
	fmt.Println("Get Clientset successful")

	//创建socket
	rest.InitRestSocket(clientSet)

	////防止主进程退出
	//ch := make(chan int )
	//select {
	//case <-ch:
	//}
}

