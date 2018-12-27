package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"net/http"
	"net/url"
	"path/filepath"
	"sync"

	"k8s.io/client-go/kubernetes"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (configFile *ConfigFile) openFile() {

	if home := homedir.HomeDir(); home != "" { // check if machine has home directory.
		// re
		// ad kubeconfig flag. if not provided use config file $HOME/.kube/config
		configFile.kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		configFile.kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// build configuration from the config file.
	config, err := clientcmd.BuildConfigFromFlags("", *configFile.kubeconfig)
	if err != nil {
		panic(err)
	}
	// create kubernetes clientset. this clientset can be used to create,delete,patch,list etc for the kubernetes resources
	configFile.clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

}

func (configFile *ConfigFile) UpdatePod(oldPodName string) (*v1.Pod, error) {

	var mutex sync.Mutex

	mutex.Lock()
	defer mutex.Unlock()

	fmt.Println("  oldPodName", oldPodName)

	pods, _ := configFile.clientset.CoreV1().Pods("default").List(metav1.ListOptions{FieldSelector: fields.SelectorFromSet(fields.Set{"metadata.name": oldPodName}).String()})
	fmt.Println(" podsssssssssss", pods)
	graceperiod := int64(0)
	modPod := pods.Items[0]
	modPod.DeletionTimestamp = nil
	err1 := configFile.clientset.CoreV1().Pods("default").Delete(oldPodName, &metav1.DeleteOptions{
		GracePeriodSeconds: &graceperiod,
	})
	if err1 != nil {
		return nil, err1
	}

	newPod := getNewPod(modPod)

	pod, err := configFile.clientset.CoreV1().Pods("default").Create(newPod)
	if err != nil {
		return nil, err
	}

	return pod, nil
}

func (configFile *ConfigFile) UpdatePod2(oldPodName string) {

	var mutex sync.Mutex

	mutex.Lock()
	defer mutex.Unlock()

	fmt.Println("  oldPodName", oldPodName)

	//pods, _ := configFile.clientset.CoreV1().Pods("default").List(metav1.ListOptions{FieldSelector: fields.SelectorFromSet(fields.Set{"metadata.name": oldPodName}).String()})
	//fmt.Println(" podsssssssssss TIME NEWWWWWWWWWWW", pods.Items[0].Name)

	// Retrieve the latest version of Deployment before attempting update
	// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
	result, getErr := configFile.clientset.CoreV1().Pods("default").Get(oldPodName, metav1.GetOptions{})
	if getErr != nil {
		panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
	}

	result.Spec.Containers[0].Image = "gofrane/cudatest" // change nginx version
	_, updateErr := configFile.clientset.CoreV1().Pods("default").Update(result)
	fmt.Println(updateErr)
}

func getNewPod(oldPod v1.Pod) *v1.Pod {
	oldContainer := oldPod.Spec.Containers[0]

	newPod := &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       oldPod.Kind,
			APIVersion: oldPod.APIVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            oldPod.GetName(),
			Namespace:       oldPod.GetNamespace(),
			Labels:          oldPod.GetLabels(),
			Annotations:     oldPod.GetAnnotations(),
			OwnerReferences: oldPod.GetOwnerReferences(),
			Initializers:    oldPod.GetInitializers(),
			Finalizers:      oldPod.GetFinalizers(),
			ClusterName:     oldPod.GetClusterName(),
		},
		Spec: v1.PodSpec{

			Volumes:        oldPod.Spec.Volumes,
			InitContainers: oldPod.Spec.InitContainers,
			Containers: []v1.Container{
				{
					Name:                     oldContainer.Name,
					Image:                    "gofrane/cudatest",
					Command:                  oldContainer.Command,
					Args:                     oldContainer.Args,
					WorkingDir:               oldContainer.WorkingDir,
					Ports:                    oldContainer.Ports,
					EnvFrom:                  oldContainer.EnvFrom,
					Env:                      oldContainer.Env,
					VolumeMounts:             oldContainer.VolumeMounts,
					LivenessProbe:            oldContainer.LivenessProbe,
					ReadinessProbe:           oldContainer.ReadinessProbe,
					Lifecycle:                oldContainer.Lifecycle,
					TerminationMessagePath:   oldContainer.TerminationMessagePath,
					TerminationMessagePolicy: oldContainer.TerminationMessagePolicy,
					ImagePullPolicy:          oldContainer.ImagePullPolicy,
					SecurityContext:          oldContainer.SecurityContext,
					Stdin:                    oldContainer.Stdin,
					StdinOnce:                oldContainer.StdinOnce,
					TTY:                      oldContainer.TTY,
				},
			},
			RestartPolicy:                 oldPod.Spec.RestartPolicy,
			TerminationGracePeriodSeconds: oldPod.Spec.TerminationGracePeriodSeconds,
			ActiveDeadlineSeconds:         oldPod.Spec.ActiveDeadlineSeconds,
			DNSPolicy:                     oldPod.Spec.DNSPolicy,
			NodeSelector:                  oldPod.Spec.NodeSelector,
			ServiceAccountName:            oldPod.Spec.ServiceAccountName,
			AutomountServiceAccountToken:  oldPod.Spec.AutomountServiceAccountToken,
			NodeName:                      oldPod.Spec.NodeName,
			HostNetwork:                   oldPod.Spec.HostNetwork,
			HostPID:                       oldPod.Spec.HostPID,
			HostIPC:                       oldPod.Spec.HostIPC,
			SecurityContext:               oldPod.Spec.SecurityContext,
			ImagePullSecrets:              oldPod.Spec.ImagePullSecrets,
			Hostname:                      oldPod.Spec.Hostname,
			Subdomain:                     oldPod.Spec.Subdomain,
			Affinity:                      oldPod.Spec.Affinity,
			SchedulerName:                 oldPod.Spec.SchedulerName,
			Tolerations:                   oldPod.Spec.Tolerations,
			HostAliases:                   oldPod.Spec.HostAliases,
		},
	}
	return newPod

}

func addnewpod() error {

	pod := Pod{
		ApiVersion: "batch/v1",
		Kind:       "Job",
		Metadata:   Metadata{Name: "pi"},
		Spec: PodSpec{
			NodeName: "",
			Containers: []Container{{
				Name:      "pi",
				Resources: ResourceRequirements{Requests: nil, Limits: nil},
				Image:     "perl",
			},
			},
		},
	}

	var b []byte
	body := bytes.NewBuffer(b)
	err := json.NewEncoder(body).Encode(pod)
	if err != nil {
		return err
	}

	request := &http.Request{
		Body:          ioutil.NopCloser(body),
		ContentLength: int64(body.Len()),
		Header:        make(http.Header),
		Method:        http.MethodPost,
		URL: &url.URL{
			Host:   apiHost,
			Path:   fmt.Sprintf(podsEndpoint, pod.Metadata.Name),
			Scheme: "http",
		},
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	if resp.StatusCode != 201 {
		return errors.New("Binding: Unexpected HTTP status code" + resp.Status)
	}
	return err

}
