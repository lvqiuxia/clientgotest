package rest

import (
	"encoding/json"
	"errors"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	"log"
	"net/http"
	"practiceCode/clientgotest/clientgo"
	"time"
)

//pod信息 json结构体
type PodMsg struct{
	Name            string
	NameSpace       string
	ReplicaSet      int32
	Image           string
	ServiceAccount  string
}

//deployment信息 json结构体
type DeploymentMsg struct{
	Name            string
	NameSpace       string
	ReplicaSet      int32
	Image           string
	ServiceAccount  string
}

//statefulSet信息 json结构体
type StatefulSetMsg struct{
	Name            string
	NameSpace       string
	ReplicaSet      int32
	Image           string
	ServiceAccount  string
}

//deamonSet信息 json结构体
type DeamonSetMsg struct{
	Name            string
	NameSpace       string
	Image           string
	ServiceAccount  string
}

//define handlerFunc
type handlerFunc struct {
	srv *http.Server
	HttpMsg *clientgo.ClientSetMsg
}

//实现handlerFunc的ServerHTTP方法
func (f handlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("msg come into")

	f.ReceiveClientMsg(w,r)
}


/* *******************************************************************
   Function:     InitRestSocket()
   Description:  初始化socket，注册处理函数，用于接收http消息
   Date:         2019.3.20
   Auther:       lvqiuxia
   Input/Output: clientset/nil
   Others:
********************************************************************** */
func InitRestSocket(msg *clientgo.ClientSetMsg){

	mux := http.NewServeMux()

	handler := handlerFunc{HttpMsg:msg}

	//handler.router = httprouter.New()
	//handler.router.GET("/sayhello",sayhello)

	//①路由注册
	mux.Handle("/", handler)
	//mux.Handlefunc("/", sayhelloName)

	s := &http.Server{
		Addr: ":8085",
		Handler: mux, //指定路由或处理器，不指定时为nil，表示使用默认的路由DefaultServeMux
		ReadTimeout: 20 * time.Second,
		WriteTimeout: 20 * time.Second,
		MaxHeaderBytes: 1 << 20,
		//ConnState: //指定连接conn的状态改变时的处理函数

	}

	//②服务监听
	err := s.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

//处理http发来的消息
func (f handlerFunc) ReceiveClientMsg (w http.ResponseWriter, r *http.Request){

	switch r.URL.Path {
	case "/createPod":
		fmt.Println("starting to create pod............")
		pod := f.CreatePod(r)
		fmt.Fprintf(w, "create pod successfull，pod name is：")
		fmt.Fprintf(w, pod)
	case "/updatePod":

	case "/listPod":
		fmt.Println("starting to list pod............")
		podList := f.ListPodMsg()
		fmt.Printf("There are %d pods in the cluster\n", len(podList))
		str := fmt.Sprintf("There are %d pods in the cluster\n", len(podList))
		fmt.Fprintf(w, str)
		for i := 0; i < len(podList); i ++ {
			fmt.Fprintf(w, podList[i].Name)
			fmt.Fprintf(w, "\n")
		}
	case "/watchPod":
		fmt.Println("starting to watch pod............")
		f.WatchPodStatus(w, r)
	case "/getPod":
		fmt.Println("starting to get pod............")
		podMsg := f.GetPodMsg(r)
		fmt.Fprintf(w, "get pod msg successfull\n")
		fmt.Fprintf(w, podMsg.ClusterName)
		fmt.Fprintf(w, podMsg.GetResourceVersion())

	case "/deletePod":
		fmt.Println("starting to delete pod............")
		err := f.DelPod(r)
		if err != nil {
			fmt.Println("delete pod failed", err)
			fmt.Fprintf(w, "delete pod failed\n")
		}else{
			fmt.Fprintf(w, "delete pod successfull\n")
		}

	case "/createDeployment":
		fmt.Println("starting to create deployment..........")
		dep := f.CreateDeployment(r)
		fmt.Fprintf(w, "create deployment successfull，deployment name is：")
		fmt.Fprintf(w, dep)

	case "/scaleDeployment":
		fmt.Println("starting to scale deployment..........")
		err := f.ScaleDeployment(r)
		if err != nil{
			fmt.Fprintf(w, "scale deployment failed ")
		}
		fmt.Fprintf(w, "scale deployment successfull ")

	case "listDeployment":
		fmt.Println("starting to list Deployment............")
		depList := f.ListDeploymentMsg(r)
		fmt.Printf("There are %d Deployments in the namespace\n", len(depList))
		str := fmt.Sprintf("There are %d Deployments in the namespace\n", len(depList))
		fmt.Fprintf(w, str)
		for i := 0; i < len(depList); i ++ {
			fmt.Fprintf(w, depList[i].Name)
			fmt.Fprintf(w, "\n")
		}

	case "/deleteDeployment":
		fmt.Println("starting to delete deployment..........")
		err := f.DelDeployment(r)
		if err != nil {
			fmt.Println("delete Deployment failed", err)
			fmt.Fprintf(w, "delete Deployment failed\n")
		}else{
			fmt.Fprintf(w, "delete Deployment successfull\n")
		}


	case "/createStatefulSet":
		fmt.Println("starting to create statefulSet..........")
		sta := f.CreateStatefulSet(r)
		fmt.Fprintf(w, "create StatefulSet successfull，statefulSet name is：")
		fmt.Fprintf(w, sta)

	case "/scaleStatefulSet":
		fmt.Println("starting to scale statefulSet..........")
		err := f.ScaleStatefulSet(r)
		if err != nil{
			fmt.Fprintf(w, "scale StatefulSet failed ")
		}
		fmt.Fprintf(w, "scale StatefulSet successfull ")

	case "/listStatefulSet":
		fmt.Println("starting to list StatefulSet............")
		staList := f.ListStatefulSetMsg(r)
		fmt.Printf("There are %d StatefulSets in the namespace\n", len(staList))
		str := fmt.Sprintf("There are %d StatefulSets in the namespace\n", len(staList))
		fmt.Fprintf(w, str)
		for i := 0; i < len(staList); i ++ {
			fmt.Fprintf(w, staList[i].Name)
			fmt.Fprintf(w, "\n")
		}

	case "/deleteStatefulSet":
		fmt.Println("starting to delete StatefulSet..........")
		err := f.DelStatefulSet(r)
		if err != nil {
			fmt.Println("delete StatefulSet failed", err)
			fmt.Fprintf(w, "delete StatefulSet failed\n")
		}else{
			fmt.Fprintf(w, "delete StatefulSet successfull\n")
		}

	case "/createDaemonSet":
		fmt.Println("starting to create DaemonSet..........")
		sta := f.CreateDeamonSet(r)
		fmt.Fprintf(w, "create DaemonSet successfull，statefulSet name is：")
		fmt.Fprintf(w, sta)

	case "/listDaemonSet":
		fmt.Println("starting to list DaemonSet............")
		deaList := f.ListDaemonSetMsg(r)
		fmt.Printf("There are %d DaemonSets in the cluster\n", len(deaList))
		str := fmt.Sprintf("There are %d DaemonSets in the cluster\n", len(deaList))
		fmt.Fprintf(w, str)
		for i := 0; i < len(deaList); i ++ {
			fmt.Fprintf(w, deaList[i].Name)
			fmt.Fprintf(w, "\n")
		}

	case "/deleteDaemonSet":
		fmt.Println("starting to delete DaemonSet..........")
		err := f.DelDaemonSet(r)
		if err != nil {
			fmt.Println("delete DaemonSet failed", err)
			fmt.Fprintf(w, "delete DaemonSet failed\n")
		}else{
			fmt.Fprintf(w, "delete DaemonSet successfull\n")
		}
	}

}

//获取所有的pod信息
func (f handlerFunc) ListPodMsg ()  []corev1.Pod{
	pods, err := f.HttpMsg.CoreV1().Pods("").List(metav1.ListOptions{})
	fmt.Println("pod msg",pods.Items[0].Name)
	if err != nil {
		panic(err.Error())
	}
	return pods.Items
}

//create pod
func (f handlerFunc) CreatePod(reqMsg *http.Request) string{
	msg := PodMsg{}
	if err := json.NewDecoder(reqMsg.Body).Decode(&msg); err != nil {
		reqMsg.Body.Close()
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
func (f handlerFunc) WatchPodStatus(w http.ResponseWriter,reqMsg *http.Request) {
	msg := PodMsg{}
	if err := json.NewDecoder(reqMsg.Body).Decode(&msg); err != nil {
		reqMsg.Body.Close()
		fmt.Println("err",err)
		log.Fatal(err)
		return
	}
	result, err := f.HttpMsg.CoreV1().Pods(msg.NameSpace).Watch(metav1.ListOptions{})
	if err != nil{
		panic(err.Error())
	}
	fmt.Fprintf(w,"watch pod in namespace successfull")

	err = WatchResult(w,result)
	fmt.Fprintf(w,"pod deleted")
	if err != nil {
		fmt.Println("pod deleted fail")
	}

	defer result.Stop()
}

//watch result
func WatchResult(w http.ResponseWriter, watcher watch.Interface) error{
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
				fmt.Fprintf(w,"pod deleted")
				//return nil
			}
		}
	}
	return nil
}

//get pod msg
func (f handlerFunc) GetPodMsg(reqMsg *http.Request) *corev1.Pod{
	msg := PodMsg{}
	if err := json.NewDecoder(reqMsg.Body).Decode(&msg); err != nil {
		reqMsg.Body.Close()
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
func (f handlerFunc) DelPod(reqMsg *http.Request) error{
	msg := PodMsg{}
	if err := json.NewDecoder(reqMsg.Body).Decode(&msg); err != nil {
		reqMsg.Body.Close()
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
func (f handlerFunc) CreateDeployment(reqMsg *http.Request) string{
	msg := DeamonSetMsg{}
	if err := json.NewDecoder(reqMsg.Body).Decode(&msg); err != nil {
		reqMsg.Body.Close()
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

func int32Ptr(i int32) *int32 { return &i }


//scale deployment
func (f handlerFunc)ScaleDeployment(reqMsg *http.Request) error{
	msg := DeploymentMsg{}
	if err := json.NewDecoder(reqMsg.Body).Decode(&msg); err != nil {
		reqMsg.Body.Close()
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
func (f handlerFunc) DelDeployment(reqMsg *http.Request) error{
	msg := DeploymentMsg{}
	if err := json.NewDecoder(reqMsg.Body).Decode(&msg); err != nil {
		reqMsg.Body.Close()
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
func(f handlerFunc) CreateStatefulSet(reqMsg *http.Request) string{
	msg := StatefulSetMsg{}
	if err := json.NewDecoder(reqMsg.Body).Decode(&msg); err != nil {
		reqMsg.Body.Close()
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
func (f handlerFunc) ScaleStatefulSet(reqMsg *http.Request) error{
	msg := StatefulSetMsg{}
	if err := json.NewDecoder(reqMsg.Body).Decode(&msg); err != nil {
		reqMsg.Body.Close()
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
func (f handlerFunc) DelStatefulSet(reqMsg *http.Request) error{
	msg := StatefulSetMsg{}
	if err := json.NewDecoder(reqMsg.Body).Decode(&msg); err != nil {
		reqMsg.Body.Close()
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
func (f handlerFunc) ListDeploymentMsg(reqMsg *http.Request) []appsv1.Deployment{
	msg := DeamonSetMsg{}
	if err := json.NewDecoder(reqMsg.Body).Decode(&msg); err != nil {
		reqMsg.Body.Close()
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
func (f handlerFunc) ListStatefulSetMsg(reqMsg *http.Request) []appsv1.StatefulSet{
	msg := StatefulSetMsg{}
	if err := json.NewDecoder(reqMsg.Body).Decode(&msg); err != nil {
		reqMsg.Body.Close()
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
func (f handlerFunc)CreateDeamonSet(reqMsg *http.Request)string{
	msg := DeamonSetMsg{}
	if err := json.NewDecoder(reqMsg.Body).Decode(&msg); err != nil {
		reqMsg.Body.Close()
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
func (f handlerFunc) ListDaemonSetMsg(reqMsg *http.Request) []appsv1.DaemonSet{
	msg := DeamonSetMsg{}
	if err := json.NewDecoder(reqMsg.Body).Decode(&msg); err != nil {
		reqMsg.Body.Close()
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
func (f handlerFunc) DelDaemonSet(reqMsg *http.Request) error{
	msg := DeamonSetMsg{}
	if err := json.NewDecoder(reqMsg.Body).Decode(&msg); err != nil {
		reqMsg.Body.Close()
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