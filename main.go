package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var kubeconfig *string
	var tmplFile *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	tmplFile = flag.String("template", "template.tmpl", "The template to render")
	flag.Parse()

	tmplRaw, err := ioutil.ReadFile(*tmplFile)
	if err != nil {
		log.Fatalf("Error reading template file %s", tmplFile)
	}
	tmpl := string(tmplRaw)
	templateParsed := template.Must(template.New("template").Parse(tmpl))

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	var b bytes.Buffer
	var oldParsedTemplate string
	var newParsedTemplate string
	for {
		services, err := clientset.CoreV1().Services("").List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		err = templateParsed.Execute(&b, services)
		newParsedTemplate = b.String()
		if err != nil {
			log.Println("Error executing template:", err)
		}
		if newParsedTemplate != oldParsedTemplate {
			oldParsedTemplate = newParsedTemplate
			log.Print(oldParsedTemplate)
		}
		b.Reset()
		time.Sleep(10 * time.Second)
	}

}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
