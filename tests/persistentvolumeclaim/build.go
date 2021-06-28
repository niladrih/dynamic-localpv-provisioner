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
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type PVCOption func(*corev1.PersistentVolumeClaim) error

func NewPVC(opts ...PVCOption) (*corev1.PersistentVolumeClaim, error) {
	pvc := &corev1.PersistentVolumeClaim{}

	var err error
	for _, opt := range opts {
		err = opt(pvc)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to build PersistentVolumeClaim.")
		}
	}

	return pvc, nil
}

func WithName(name string) PVCOption {
	return func(pvc *corev1.PersistentVolumeClaim) error {
		if len(name) == 0 {
			return errors.New("Failed to set Name. Name is an empty string.")
		}

		pvc.ObjectMeta.Name = name
		return nil
	}
}

func WithGenerateName(generateName string) PVCOption {
	return func(pvc *corev1.PersistentVolumeClaim) error {
		if len(generateName) == 0 {
			return errors.New("Failed to set GenerateName. Name prefix is an empty string.")
		}

		pvc.ObjectMeta.GenerateName = generateName + "-"
		return nil
	}
}

func WithNamespace(namespace string) PVCOption {
	return func(pvc *corev1.PersistentVolumeClaim) error {
		if len(namespace) == 0 {
			namespace = "default"
		}

		pvc.ObjectMeta.Namespace = namespace
		return nil
	}
}

func WithStorageClass(storageClass string) PVCOption {
	return func(pvc *corev1.PersistentVolumeClaim) error {
		if len(storageClass) == 0 {
			return errors.New("Failed to set StorageClassName." +
				" StorageClassName is an empty string.")
		}

		pvc.Spec.StorageClassName = &storageClass
		return nil
	}
}

func WithAccessModes(accessModes ...corev1.PersistentVolumeAccessMode) PVCOption {
	return func(pvc *corev1.PersistentVolumeClaim) error {
		if len(accessModes) == 0 {
			return errors.New("Failed to set AccessMode(s). Too few arguments.")
		}

		pvc.Spec.AccessModes = accessModes
		return nil
	}
}

func WithCapacity(capacity string) PVCOption {
	return func(pvc *corev1.PersistentVolumeClaim) error {
		resCapacity, err := resource.ParseQuantity(capacity)
		if err != nil {
			return errors.New("Failed to build PVC. Failed to parse capacity " +
				capacity)
		}

		resourceList := corev1.ResourceList{
			corev1.ResourceName(corev1.ResourceStorage): resCapacity,
		}

		pvc.Spec.Resources.Requests = resourceList
		return nil
	}
}

func WithVolumeMode(volumeMode corev1.PersistentVolumeMode) PVCOption {
	return func(pvc *corev1.PersistentVolumeClaim) error {
		if len(volumeMode) == 0 {
			volumeMode = "Filesystem"
		}

		pvc.Spec.VolumeMode = &volumeMode
		return nil
	}
}
