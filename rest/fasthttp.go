package rest

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/valyala/fasthttp"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"log"
	"practiceCode/clientgotest/clientgo"
)

////pod信息 json结构体
//type PodMsg struct{
//	Name            string
//	NameSpace       string
//	ReplicaSet      int32
//	Image           string
//	ServiceAccount  string
//}
//
////deployment信息 json结构体
//type DeploymentMsg struct{
//	Name            string
//	NameSpace       string
//	ReplicaSet      int32
//	Image           string
//	ServiceAccount  string
//}
//
////statefulSet信息 json结构体
//type StatefulSetMsg struct{
//	Name            string
//	NameSpace       string
//	ReplicaSet      int32
//	Image           string
//	ServiceAccount  string
//}
//
////deamonSet信息 json结构体
//type DeamonSetMsg struct{
//	Name            string
//	NameSpace       string
//	Image           string
//	ServiceAccount  string
//}

var addr = flag.String("addr", ":8085", "TCP address to listen to")

//define handlerFunc
type fastHandlerFunc struct {
	HttpMsg *clientgo.ClientSetMsg
}

//处理函数
func (f fastHandlerFunc) HandleFastHTTP(ctx *fasthttp.RequestCtx){
	f.ReceiveFastHttpMsg(ctx)
}


/* *******************************************************************
   Function:     InitFastHttpSocket()
   Description:  初始化fasthttp，注册处理函数，用于接收http消息
   Date:         2019.3.25
   Auther:       lvqiuxia
   Input/Output: clientset/nil
   Others:
********************************************************************** */
func InitFastHttpSocket(msg *clientgo.ClientSetMsg){

	flag.Parse()

	handler := &fastHandlerFunc{HttpMsg:msg}

	//h = fasthttp.CompressHandler(h)

	//②服务监听
	if err := fasthttp.ListenAndServe(*addr, handler.HandleFastHTTP); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

//处理http发来的消息
func (f fastHandlerFunc) ReceiveFastHttpMsg (ctx *fasthttp.RequestCtx){

	switch string(ctx.Path()) {
	case "/createPod":
		fmt.Println("starting to create pod............")
		pod := f.CreatePod(ctx)
		fmt.Fprintf(ctx, "create pod successfull，pod name is：")
		fmt.Fprintf(ctx, pod)
	case "/updatePod":

	case "/listPod":
		fmt.Println("starting to list pod............")
		podList := f.ListPodMsg()
		fmt.Printf("There are %d pods in the cluster\n", len(podList))
		str := fmt.Sprintf("There are %d pods in the cluster\n", len(podList))
		fmt.Fprintf(ctx, str)
		for i := 0; i < len(podList); i ++ {
			fmt.Fprintf(ctx, podList[i].Name)
			fmt.Fprintf(ctx, "\n")
		}
	case "/watchPod":
		fmt.Println("starting to watch pod............")
		f.WatchPodStatus(ctx)
	case "/getPod":
		fmt.Println("starting to get pod............")
		podMsg := f.GetPodMsg(ctx)
		fmt.Fprintf(ctx, "get pod msg successfull\n")
		fmt.Fprintf(ctx, podMsg.ClusterName)
		fmt.Fprintf(ctx, podMsg.GetResourceVersion())

	case "/deletePod":
		fmt.Println("starting to delete pod............")
		err := f.DelPod(ctx)
		if err != nil {
			fmt.Println("delete pod failed", err)
			fmt.Fprintf(ctx, "delete pod failed\n")
		}else{
			fmt.Fprintf(ctx, "delete pod successfull\n")
		}

	case "/createDeployment":
		fmt.Println("starting to create deployment..........")
		dep := f.CreateDeployment(ctx)
		fmt.Fprintf(ctx, "create deployment successfull，deployment name is：")
		fmt.Fprintf(ctx, dep)

	case "/scaleDeployment":
		fmt.Println("starting to scale deployment..........")
		err := f.ScaleDeployment(ctx)
		if err != nil{
			fmt.Fprintf(ctx, "scale deployment failed ")
		}
		fmt.Fprintf(ctx, "scale deployment successfull ")

	case "listDeployment":
		fmt.Println("starting to list Deployment............")
		depList := f.ListDeploymentMsg(ctx)
		fmt.Printf("There are %d Deployments in the namespace\n", len(depList))
		str := fmt.Sprintf("There are %d Deployments in the namespace\n", len(depList))
		fmt.Fprintf(ctx, str)
		for i := 0; i < len(depList); i ++ {
			fmt.Fprintf(ctx, depList[i].Name)
			fmt.Fprintf(ctx, "\n")
		}

	case "/deleteDeployment":
		fmt.Println("starting to delete deployment..........")
		err := f.DelDeployment(ctx)
		if err != nil {
			fmt.Println("delete Deployment failed", err)
			fmt.Fprintf(ctx, "delete Deployment failed\n")
		}else{
			fmt.Fprintf(ctx, "delete Deployment successfull\n")
		}


	case "/createStatefulSet":
		fmt.Println("starting to create statefulSet..........")
		sta := f.CreateStatefulSet(ctx)
		fmt.Fprintf(ctx, "create StatefulSet successfull，statefulSet name is：")
		fmt.Fprintf(ctx, sta)

	case "/scaleStatefulSet":
		fmt.Println("starting to scale statefulSet..........")
		err := f.ScaleStatefulSet(ctx)
		if err != nil{
			fmt.Fprintf(ctx, "scale StatefulSet failed ")
		}
		fmt.Fprintf(ctx, "scale StatefulSet successfull ")

	case "/listStatefulSet":
		fmt.Println("starting to list StatefulSet............")
		staList := f.ListStatefulSetMsg(ctx)
		fmt.Printf("There are %d StatefulSets in the namespace\n", len(staList))
		str := fmt.Sprintf("There are %d StatefulSets in the namespace\n", len(staList))
		fmt.Fprintf(ctx, str)
		for i := 0; i < len(staList); i ++ {
			fmt.Fprintf(ctx, staList[i].Name)
			fmt.Fprintf(ctx, "\n")
		}

	case "/deleteStatefulSet":
		fmt.Println("starting to delete StatefulSet..........")
		err := f.DelStatefulSet(ctx)
		if err != nil {
			fmt.Println("delete StatefulSet failed", err)
			fmt.Fprintf(ctx, "delete StatefulSet failed\n")
		}else{
			fmt.Fprintf(ctx, "delete StatefulSet successfull\n")
		}

	case "/createDaemonSet":
		fmt.Println("starting to create DaemonSet..........")
		sta := f.CreateDeamonSet(ctx)
		fmt.Fprintf(ctx, "create DaemonSet successfull，statefulSet name is：")
		fmt.Fprintf(ctx, sta)

	case "/listDaemonSet":
		fmt.Println("starting to list DaemonSet............")
		deaList := f.ListDaemonSetMsg(ctx)
		fmt.Printf("There are %d DaemonSets in the cluster\n", len(deaList))
		str := fmt.Sprintf("There are %d DaemonSets in the cluster\n", len(deaList))
		fmt.Fprintf(ctx, str)
		for i := 0; i < len(deaList); i ++ {
			fmt.Fprintf(ctx, deaList[i].Name)
			fmt.Fprintf(ctx, "\n")
		}

	case "/deleteDaemonSet":
		fmt.Println("starting to delete DaemonSet..........")
		err := f.DelDaemonSet(ctx)
		if err != nil {
			fmt.Println("delete DaemonSet failed", err)
			fmt.Fprintf(ctx, "delete DaemonSet failed\n")
		}else{
			fmt.Fprintf(ctx, "delete DaemonSet successfull\n")
		}
	}

}

//获取所有的pod信息
func (f fastHandlerFunc) ListPodMsg ()  []corev1.Pod{
	pods, err := f.HttpMsg.CoreV1().Pods("").List(metav1.ListOptions{})
	fmt.Println("pod msg",pods.Items[0].Name)
	if err != nil {
		panic(err.Error())
	}
	return pods.Items
}

//create pod
func (f fastHandlerFunc) CreatePod(reqMsg *fasthttp.RequestCtx) string{
	msg := PodMsg{}
	if err := json.Unmarshal(reqMsg.Request.Body(),&msg); err != nil {
		fmt.Println("err",err)
		log.Fatal(err)
		return ""
	}
	pod := &corev1.Pod{
		ObjectMeta:metav1.ObjectMeta{Name:msg.Name},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "web",
					Image: msg.Image,
					Ports: []corev1.ContainerPort{
						{
							Name:          "http",
							Protocol:      corev1.ProtocolTCP,
							ContainerPort: 80,
						},
					},
				},
			},
			ServiceAccountName:msg.ServiceAccount,
		},
	}
	result, err:= f.HttpMsg.CoreV1().Pods(msg.NameSpace).Create(pod)
	if err != nil{
		panic(err.Error())
	}
	fmt.Printf("Created pod %q in namespace %q.\n", result.GetObjectMeta().GetName(),result.GetObjectMeta().GetNamespace())
	return result.GetObjectMeta().GetName()
}

//watch pod
func (f fastHandlerFunc) WatchPodStatus(reqMsg *fasthttp.RequestCtx) {
	msg := PodMsg{}
	if err := json.Unmarshal(reqMsg.Request.Body(),&msg); err != nil {
		fmt.Println("err",err)
		log.Fatal(err)
		return
	}
	result, err := f.HttpMsg.CoreV1().Pods(msg.NameSpace).Watch(metav1.ListOptions{})
	if err != nil{
		panic(err.Error())
	}
	fmt.Fprintf(reqMsg,"watch pod in namespace successfull")

	err = f.WatchResult(reqMsg,result)
	fmt.Fprintf(reqMsg,"pod deleted")
	if err != nil {
		fmt.Println("pod deleted fail")
	}

	defer result.Stop()
}

//watch result
func (f fastHandlerFunc) WatchResult(reqMsg *fasthttp.RequestCtx, watcher watch.Interface) error{
	fmt.Println("come into WatchResult")
	ch := watcher.ResultChan()

	for {
		select {
		case event, ok := <-ch:
			if !ok {
				return errors.New("error")
			}
			fmt.Println("something happen to the pod")
			// We need to inspect the event and get ResourceVersion out of it
			switch event.Type {
			case watch.Added:
			case watch.Modified:
			case watch.Deleted:
				fmt.Println("pod deleted")
				fmt.Fprintf(reqMsg,"pod deleted")
				//return nil
			}
		}
	}
	return nil
}

//get pod msg
func (f fastHandlerFunc) GetPodMsg(reqMsg *fasthttp.RequestCtx) *corev1.Pod{
	msg := PodMsg{}
	if err := json.Unmarshal(reqMsg.Request.Body(),&msg); err != nil {
		fmt.Println("err",err)
		log.Fatal(err)
		return nil
	}
	pods, err := f.HttpMsg.CoreV1().Pods(msg.NameSpace).Get(msg.Name,metav1.GetOptions{})

	if err != nil {
		panic(err.Error())
	}

	return pods
}

//delete pod
func (f fastHandlerFunc) DelPod(reqMsg *fasthttp.RequestCtx) error{
	msg := PodMsg{}
	if err := json.Unmarshal(reqMsg.Request.Body(),&msg); err != nil {
		fmt.Println("err",err)
		log.Fatal(err)
		return err
	}
	err := f.HttpMsg.CoreV1().Pods(msg.NameSpace).Delete(msg.Name,&metav1.DeleteOptions{})

	if err != nil {
		return err
	}

	return nil
}


//create deploymenrt
func (f fastHandlerFunc) CreateDeployment(reqMsg *fasthttp.RequestCtx) string{
	msg := DeamonSetMsg{}
	if err := json.Unmarshal(reqMsg.Request.Body(),&msg); err != nil {
		fmt.Println("err",err)
		log.Fatal(err)
		return ""
	}
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: msg.Name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "webtest",
							Image: msg.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	result, err:= f.HttpMsg.AppsV1().Deployments(msg.NameSpace).Create(deployment)
	if err != nil{
		panic(err.Error())
	}
	fmt.Printf("Created deployment %q in namespace %q.\n", result.GetObjectMeta().GetName(),
		result.GetObjectMeta().GetNamespace())
	return result.GetObjectMeta().GetName()
}



//scale deployment
func (f fastHandlerFunc)ScaleDeployment(reqMsg *fasthttp.RequestCtx) error{
	msg := DeploymentMsg{}
	if err := json.Unmarshal(reqMsg.Request.Body(),&msg); err != nil {
		fmt.Println("err",err)
		log.Fatal(err)
		return err
	}
	scale := &autoscalingv1.Scale{
		ObjectMeta: metav1.ObjectMeta{
			Name: msg.Name,
			Namespace:msg.NameSpace,
		},
		Spec: autoscalingv1.ScaleSpec{msg.ReplicaSet},
	}
	_, err := f.HttpMsg.AppsV1().Deployments(msg.NameSpace).UpdateScale(msg.Name,scale)
	if err != nil{
		panic(err.Error())
		return err
	}
	return nil
}

//delete deployment
func (f fastHandlerFunc) DelDeployment(reqMsg *fasthttp.RequestCtx) error{
	msg := DeploymentMsg{}
	if err := json.Unmarshal(reqMsg.Request.Body(),&msg); err != nil {
		fmt.Println("err",err)
		log.Fatal(err)
		return err
	}
	err := f.HttpMsg.AppsV1().Deployments(msg.NameSpace).Delete(msg.Name,&metav1.DeleteOptions{})

	if err != nil {
		return err
	}

	return nil
}


//create statefulSet
func(f fastHandlerFunc) CreateStatefulSet(reqMsg *fasthttp.RequestCtx) string{
	msg := StatefulSetMsg{}
	if err := json.Unmarshal(reqMsg.Request.Body(),&msg); err != nil {
		fmt.Println("err",err)
		log.Fatal(err)
		return ""
	}

	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: msg.Name,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "webtest",
							Image: msg.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	result, err := f.HttpMsg.AppsV1().StatefulSets(msg.NameSpace).Create(statefulSet)
	if err != nil{
		panic(err.Error())
	}
	fmt.Printf("Created StatefulSet %q in namespace %q.\n", result.GetObjectMeta().GetName(),
		result.GetObjectMeta().GetNamespace())
	return result.GetObjectMeta().GetNamespace()
}

//scale statefulSet
func (f fastHandlerFunc) ScaleStatefulSet(reqMsg *fasthttp.RequestCtx) error{
	msg := StatefulSetMsg{}
	if err := json.Unmarshal(reqMsg.Request.Body(),&msg); err != nil {
		fmt.Println("err",err)
		log.Fatal(err)
		return err
	}

	scale := &autoscalingv1.Scale{
		ObjectMeta: metav1.ObjectMeta{
			Name: msg.Name,
			Namespace:msg.NameSpace,
		},
		Spec: autoscalingv1.ScaleSpec{msg.ReplicaSet},
	}

	f.HttpMsg.AppsV1().StatefulSets(msg.NameSpace).UpdateScale(msg.Name,scale)

	return nil
}


//delete StatefulSets
func (f fastHandlerFunc) DelStatefulSet(reqMsg *fasthttp.RequestCtx) error{
	msg := StatefulSetMsg{}
	if err := json.Unmarshal(reqMsg.Request.Body(),&msg); err != nil {
		fmt.Println("err",err)
		log.Fatal(err)
		return err
	}
	err := f.HttpMsg.AppsV1().StatefulSets(msg.NameSpace).Delete(msg.Name,&metav1.DeleteOptions{})

	if err != nil {
		return err
	}

	return nil
}



//list deployment
func (f fastHandlerFunc) ListDeploymentMsg(reqMsg *fasthttp.RequestCtx) []appsv1.Deployment{
	msg := DeamonSetMsg{}
	if err := json.Unmarshal(reqMsg.Request.Body(),&msg); err != nil {
		fmt.Println("err",err)
		log.Fatal(err)
		return nil
	}
	depMsg,err := f.HttpMsg.AppsV1().Deployments(msg.NameSpace).List(metav1.ListOptions{})
	if err != nil{
		fmt.Println("list deployment failed",err)
		return nil
	}

	return depMsg.Items
}


//list StatefulSet
func (f fastHandlerFunc) ListStatefulSetMsg(reqMsg *fasthttp.RequestCtx) []appsv1.StatefulSet{
	msg := StatefulSetMsg{}
	if err := json.Unmarshal(reqMsg.Request.Body(),&msg); err != nil {
		fmt.Println("err",err)
		log.Fatal(err)
		return nil
	}
	depMsg,err := f.HttpMsg.AppsV1().StatefulSets(msg.NameSpace).List(metav1.ListOptions{})
	if err != nil{
		fmt.Println("list StatefulSet failed",err)
		return nil
	}

	return depMsg.Items
}


//create deamonset
func (f fastHandlerFunc)CreateDeamonSet(reqMsg *fasthttp.RequestCtx)string{
	msg := DeamonSetMsg{}
	if err := json.Unmarshal(reqMsg.Request.Body(),&msg); err != nil {
		fmt.Println("err",err)
		log.Fatal(err)
		return ""
	}

	deamonSet := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: msg.Name,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "webtest",
							Image: msg.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	result, err := f.HttpMsg.AppsV1().DaemonSets(msg.NameSpace).Create(deamonSet)
	if err != nil{
		panic(err.Error())
	}
	fmt.Printf("Created DaemonSets %q in namespace %q.\n", result.GetObjectMeta().GetName(),
		result.GetObjectMeta().GetNamespace())
	return result.GetObjectMeta().GetNamespace()
}


//list DaemonSets
func (f fastHandlerFunc) ListDaemonSetMsg(reqMsg *fasthttp.RequestCtx) []appsv1.DaemonSet{
	msg := DeamonSetMsg{}
	if err := json.Unmarshal(reqMsg.Request.Body(),&msg); err != nil {
		fmt.Println("err",err)
		log.Fatal(err)
		return nil
	}
	depMsg,err := f.HttpMsg.AppsV1().DaemonSets(msg.NameSpace).List(metav1.ListOptions{})
	if err != nil{
		fmt.Println("list DaemonSet failed",err)
		return nil
	}

	return depMsg.Items
}


//delete DaemonSet
func (f fastHandlerFunc) DelDaemonSet(reqMsg *fasthttp.RequestCtx) error{
	msg := DeamonSetMsg{}

	if err := json.Unmarshal(reqMsg.Request.Body(),&msg); err != nil {
		fmt.Println("err",err)
		log.Fatal(err)
		return err
	}
	err := f.HttpMsg.AppsV1().DaemonSets(msg.NameSpace).Delete(msg.Name,&metav1.DeleteOptions{})

	if err != nil {
		return err
	}

	return nil
}