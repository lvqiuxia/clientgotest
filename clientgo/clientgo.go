package clientgo

import (
	"flag"
	"fmt"
	"k8s.io/client-go/discovery"
	admissionregistrationv1beta1 "k8s.io/client-go/kubernetes/typed/admissionregistration/v1beta1"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	batchv1 "k8s.io/client-go/kubernetes/typed/batch/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/flowcontrol"
	"os"
	"path/filepath"
)


//定义是以集群内还是集群外的形式获取资源对象
var (
	OutOfClusterClientConfig = "OutOfClusterClientConfig"
	InOfClusterClientConfig = "InOfClusterClientConfig"
)


//定义clientSet的结构体
type ClientSetMsg struct{
	*discovery.DiscoveryClient
	admissionregistrationV1beta1 *admissionregistrationv1beta1.AdmissionregistrationV1beta1Client
	appsV1                       *appsv1.AppsV1Client
	batchV1                      *batchv1.BatchV1Client
	coreV1                       *corev1.CoreV1Client
}



/* *******************************************************************
   Function:     CreateClientSet()
   Description:  对k8s集群的客户端进行认证，并获取到clientset，以便后续进行资源获取
   Date:         2019.3.20
   Auther:       lvqiuxia
   Input/Output: nil/clientset
   Others:
********************************************************************** */
func CreateClientSet(mode string)  (*ClientSetMsg, error){

	//集群内和集群外两种获取config的方式
	var config *restclient.Config
	switch mode{
	case OutOfClusterClientConfig:
		config = GetConfigWithOutCluster()
	case InOfClusterClientConfig:
		config = GetConfigWithInCluster()
	}

	// create the clientset
	clientset, err := NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return clientset,nil
}


//集群外的方式获取config
func GetConfigWithOutCluster()*restclient.Config{
	fmt.Println("come into GetConfigWithOutCluster")

	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute" +
			" path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	return config
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

//集群内的方式获取config
func GetConfigWithInCluster()*restclient.Config{
	fmt.Println("come into GetConfigWithInCluster")

	// creates the in-cluster config
	config, err := restclient.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	return config
}



// NewForConfig creates a new Clientset for the given config.
func NewForConfig(c *restclient.Config) (*ClientSetMsg, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}
	var cs ClientSetMsg
	var err error
	cs.admissionregistrationV1beta1, err = admissionregistrationv1beta1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.appsV1, err = appsv1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	cs.batchV1, err = batchv1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	cs.coreV1, err = corev1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	cs.DiscoveryClient, err = discovery.NewDiscoveryClientForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	return &cs, nil
}

// AdmissionregistrationV1beta1 retrieves the AdmissionregistrationV1beta1Client
func (c *ClientSetMsg) AdmissionregistrationV1beta1() admissionregistrationv1beta1.AdmissionregistrationV1beta1Interface {
	return c.admissionregistrationV1beta1
}

// AppsV1 retrieves the AppsV1Client
func (c *ClientSetMsg) AppsV1() appsv1.AppsV1Interface {
	return c.appsV1
}

// BatchV1 retrieves the BatchV1Client
func (c *ClientSetMsg) BatchV1() batchv1.BatchV1Interface {
	return c.batchV1
}

// CoreV1 retrieves the CoreV1Client
func (c *ClientSetMsg) CoreV1() corev1.CoreV1Interface {
	return c.coreV1
}

// Discovery retrieves the DiscoveryClient
func (c *ClientSetMsg) Discovery() discovery.DiscoveryInterface {
	if c == nil {
		return nil
	}
	return c.DiscoveryClient
}
