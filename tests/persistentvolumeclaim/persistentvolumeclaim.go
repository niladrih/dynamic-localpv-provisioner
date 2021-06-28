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

package persistentvolumeclaim

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	coreclientv1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

const (
	integrationTestLabelKey   = "openebs.io/integration-test"
	integrationTestLabelValue = "true"
)

func CreateForTest(pvc *corev1.PersistentVolumeClaim, client *coreclientv1.CoreV1Client) (*corev1.PersistentVolumeClaim, error) {
	pvc.ObjectMeta.Labels[integrationTestLabelKey] = integrationTestLabelValue
	return client.PersistentVolumeClaims(pvc.ObjectMeta.Namespace).Create(pvc)
}

func DeleteForTest(pvc *corev1.PersistentVolumeClaim, client *coreclientv1.CoreV1Client) error {
	return client.PersistentVolumeClaims(pvc.ObjectMeta.Namespace).Delete(pvc.Name, &metav1.DeleteOptions{})
}

func ListForTest(namespace string, client *coreclientv1.CoreV1Client) (*corev1.PersistentVolumeClaimList, error) {
	return client.PersistentVolumeClaims(namespace).List(metav1.ListOptions{
		LabelSelector: integrationTestLabelKey + "=" + integrationTestLabelValue,
	})
}
