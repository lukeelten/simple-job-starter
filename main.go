package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var config *rest.Config
var client *kubernetes.Clientset

var namespace string

func main() {
	var err error
	config, err = rest.InClusterConfig()
	if err != nil {
		panic(err)
	}

	client, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	namespaceBytes, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		panic(err)
	}

	namespace = string(namespaceBytes)

	r := gin.Default()
	r.GET("/start", func(c *gin.Context) {
		job := startJob()
		c.JSON(200, job)
	})

	r.GET("/status", func(c *gin.Context) {
		list, err := client.BatchV1().Jobs(namespace).List(metav1.ListOptions{
			LabelSelector: "app=bash-job",
		})

		if err != nil {
			panic(err)
		}

		c.JSON(200, list)
	})

	r.Run()
}

func startJob() *v1.Job {
	id := generateUUID()
	now := time.Now()
	name := fmt.Sprintf("bash-job-%d-%d-%d-%d-%d-%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	labels := make(map[string]string, 0)
	labels["app"] = "bash-job"
	labels["uuid"] = id

	job := v1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
		Spec: v1.JobSpec{
			Parallelism:  intp(1),
			Completions:  intp(1),
			BackoffLimit: intp(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
					Name:   "bash-job",
				},
				Spec: corev1.PodSpec{
					RestartPolicy:      corev1.RestartPolicyNever,
					ServiceAccountName: "bash-job",
					Containers: []corev1.Container{{
						Image:           "bash-job",
						Name:            "main",
						ImagePullPolicy: corev1.PullAlways,
					}},
				},
			},
		},
	}

	jobRet, err := client.BatchV1().Jobs(namespace).Create(&job)
	if err != nil {
		panic(err)
	}

	return jobRet
}

func intp(i int32) *int32 {
	return &i
}

func generateUUID() string {
	uid := uuid.New()
	return uid.String()
}
