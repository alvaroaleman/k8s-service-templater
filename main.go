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

	"gopkg.in/yaml.v2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var kubeconfig *string
	var configPath *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	configPath = flag.String("config", "config.yml", "The application confid")
	flag.Parse()

	appConfigRaw, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("Error reading config file: '%s'", err)
	}

	var appConfig AppConfig
	err = yaml.Unmarshal(appConfigRaw, &appConfig)
	if err != nil {
		log.Fatalf("Errir parsing config file: '%s'", err)
	}

	tmplRaw, err := ioutil.ReadFile(appConfig.Template)
	if err != nil {
		log.Fatalf("Error reading template file %s", appConfig.Template)
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

	err = run(appConfig, clientset, templateParsed)
	if err != nil {
		log.Fatalf("Error running app: '%s'", err)
	}

}

func run(appConfig AppConfig, clientset *kubernetes.Clientset, template *template.Template) error {
	var b bytes.Buffer
	var oldParsedTemplate string
	var newParsedTemplate string

	for {
		services, err := clientset.CoreV1().Services("").List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		err = template.Execute(&b, services)
		newParsedTemplate = b.String()
		if err != nil {
			log.Println("Error executing template:", err)
		}
		if newParsedTemplate != oldParsedTemplate {
			oldParsedTemplate = newParsedTemplate
			err = ioutil.WriteFile(appConfig.TemplateDestination, []byte(oldParsedTemplate), 0644)
			if err != nil {
				log.Println("Error writing template: '%s'")
			} else {
				log.Print("Wrote new config..")
			}
		}
		b.Reset()
		time.Sleep(10 * time.Second)
	}

	return nil

}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
