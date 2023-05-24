package main

import (
	"os"
)

// func main() {
// 	var kubeconfig *string
// 	if home := homeDir(); home != "" {
// 		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
// 	} else {
// 		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
// 	}
// 	flag.Parse()
// 	// uses the current context in kubeconfig
// 	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	// creates the clientset
// 	clientset, err := kubernetes.NewForConfig(config)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	for {
// 		pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
// 		time.Sleep(10 * time.Second)
// 	}
// }

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
