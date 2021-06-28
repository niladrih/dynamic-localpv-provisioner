/*
Copyright 2019 The OpenEBS Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package storageclass

import (
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	storageclientv1 "k8s.io/client-go/kubernetes/typed/storage/v1"
)

const (
	integrationTestLabelKey   = "openebs.io/integration-test"
	integrationTestLabelValue = "true"
)

func CreateForTest(sc *storagev1.StorageClass, client *storageclientv1.StorageV1Client) (*storagev1.StorageClass, error) {
	sc.ObjectMeta.Labels[integrationTestLabelKey] = integrationTestLabelValue
	return client.StorageClasses().Create(sc)
}

func DeleteForTest(sc *storagev1.StorageClass, client *storageclientv1.StorageV1Client) error {
	return client.StorageClasses().Delete(sc.Name, &metav1.DeleteOptions{})
}

func ListForTest(client *storageclientv1.StorageV1Client) (*storagev1.StorageClassList, error) {
	return client.StorageClasses().List(metav1.ListOptions{
		LabelSelector: integrationTestLabelKey + "=" + integrationTestLabelValue,
	})
}
