package main

import (
	"encoding/json"
	"net/http"
	"net/url"
)











//Get the Pod Succeeded
func getSucceededPods() (*PodList, error) {
	var PodSucceededList PodList

	v := url.Values{}

	v.Add("fieldSelector", "status.phase=Succeeded")

	request := &http.Request{
		Header: make(http.Header),
		Method: http.MethodGet,
		URL: &url.URL{
			Host:     apiHost,
			Path:     podsEndpoint,
			RawQuery: v.Encode(),
			Scheme:   "http",
		},
	}
	request.Header.Set("Accept", "application/json, */*")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(resp.Body).Decode(&PodSucceededList)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return &PodSucceededList, nil
}

func getFailedPods() (*PodList, error) {
	var PodFailedList PodList

	v := url.Values{}

	v.Add("fieldSelector", "status.phase=Failed")

	request := &http.Request{
		Header: make(http.Header),
		Method: http.MethodGet,
		URL: &url.URL{
			Host:     apiHost,
			Path:     podsEndpoint,
			RawQuery: v.Encode(),
			Scheme:   "http",
		},
	}
	request.Header.Set("Accept", "application/json, */*")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(resp.Body).Decode(&PodFailedList)
	if err != nil {
		return nil, err
	}
  defer resp.Body.Close()
	return &PodFailedList, nil
}

func getRunningPods() (*PodList, error) {
	var PodRunningList PodList

	v := url.Values{}

	v.Add("fieldSelector", "status.phase=Running")

	request := &http.Request{
		Header: make(http.Header),
		Method: http.MethodGet,
		URL: &url.URL{
			Host:     apiHost,
			Path:     podsEndpoint,
			RawQuery: v.Encode(),
			Scheme:   "http",
		},
	}
	request.Header.Set("Accept", "application/json, */*")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(resp.Body).Decode(&PodRunningList)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return &PodRunningList, nil
}

func getPendingPods() (*PodList, error) {
	var PodPendingList PodList

	v := url.Values{}

	v.Add("fieldSelector", "status.phase=Pending")

	request := &http.Request{
		Header: make(http.Header),
		Method: http.MethodGet,
		URL: &url.URL{
			Host:     apiHost,
			Path:     podsEndpoint,
			RawQuery: v.Encode(),
			Scheme:   "http",
		},
	}
	request.Header.Set("Accept", "application/json, */*")



	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(resp.Body).Decode(&PodPendingList)
	if err != nil {
		return nil, err
	}






	defer resp.Body.Close()
	return &PodPendingList, nil
}





func getImage(pod *Pod) string {

	var imagePending string

	imagePending = pod.Spec.Containers[0].Image

	return imagePending
}

