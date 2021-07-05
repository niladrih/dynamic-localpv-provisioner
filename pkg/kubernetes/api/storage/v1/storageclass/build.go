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
	mconfig "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
)

const (
	localPVcasTypeValue = "local"

	// Provisioner Name
	localPVprovisionerName = "openebs.io/local"

	// The following are imported from mconfig at the moment
	// CASConfigKey = "cas.openebs.io/config"
	// CASTypeKey = "openebs.io/cas-type"
)

type StorageClassOption func(*storagev1.StorageClass) error

func NewStorageClass(opts ...StorageClassOption) (*storagev1.StorageClass, error) {
	s := &storagev1.StorageClass{}

	var err error
	for _, opt := range opts {
		err = opt(s)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to build StorageClass.")
		}
	}

	return s, nil
}

func WithName(name string) StorageClassOption {
	return func(s *storagev1.StorageClass) error {
		if len(name) == 0 {
			return errors.New("Failed to set Name. Name is an empty string.")
		}

		s.ObjectMeta.Name = name
		return nil
	}
}

func WithGenerateName(generateName string) StorageClassOption {
	return func(s *storagev1.StorageClass) error {
		if len(generateName) == 0 {
			return errors.New("Failed to set GenerateName. Name prefix is an empty string.")
		}

		s.ObjectMeta.GenerateName = generateName + "-"
		return nil
	}
}

func WithLocalPV() StorageClassOption {
	return func(s *storagev1.StorageClass) error {
		if _, ok := s.ObjectMeta.Annotations[string(mconfig.CASTypeKey)]; ok {
			return errors.New("Annotation '" + string(mconfig.CASTypeKey) +
				"' is already set.")
		}
		if len(s.Provisioner) > 0 {
			return errors.New("Provisioner name is already set.")
		}

		// Set the cas-type annotation
		if s.ObjectMeta.Annotations == nil {
			s.ObjectMeta.Annotations = map[string]string{}
		}
		s.ObjectMeta.Annotations[string(mconfig.CASTypeKey)] = localPVcasTypeValue
		// Set the provisioner value for
		// openebs-localpv-provisioner PV controller
		s.Provisioner = localPVprovisionerName

		return nil
	}
}

func WithHostpath(hostpathDir string) StorageClassOption {
	return func(s *storagev1.StorageClass) error {
		// Check if the path is a valid one
		if !isValidPath(hostpathDir) {
			return errors.New("Invalid hostpath directory. Path" +
				" must be an absolute path and must be a " +
				"directory which is not directly under '/'.")
		}
		// Check for existing CAS config and Provisioner name
		// Check if the existing parameters are usable
		// with "hostpath" StorageType
		if !isCompatibleWithHostpath(s) {
			return errors.New("Failed to set StorageType and BasePath for Hostpath. " +
				"Invalid existing '" + string(mconfig.CASConfigKey) + "' annotation" +
				" parameters or Provisioner name.")
		}

		config := "- name: StorageType\n" +
			"  value: \"hostpath\"\n" +
			"- name: BasePath\n" +
			"  value: \"" + hostpathDir + "\"\n"

		ok := writeOrAppendCASConfig(s, config)
		if !ok {
			return errors.New("Failed to set StorageType and" +
				" BasePath parameters for Hostpath.")
		}
		return nil
	}
}

func WithDevice() StorageClassOption {
	return func(s *storagev1.StorageClass) error {
		// Check for existing CAS config and Provisioner name
		// Check if the existing parameters are usable
		// with "device" StorageType
		if !isCompatibleWithDevice(s) {
			return errors.New("Failed to set StorageType for Device. " +
				"Invalid existing '" + string(mconfig.CASConfigKey) +
				"' annotaion parameters or Provisioner name.")
		}

		config := "- name: StorageType\n" +
			"  value: \"device\"\n"

		ok := writeOrAppendCASConfig(s, config)
		if !ok {
			return errors.New("Failed to set StorageType parameter for Device.")
		}

		return nil
	}
}

func WithVolumeBindingMode(volBindingMode storagev1.VolumeBindingMode) StorageClassOption {
	return func(s *storagev1.StorageClass) error {
		if len(volBindingMode) == 0 {
			volBindingMode = "WaitForFirstConsumer"
		}

		s.VolumeBindingMode = &volBindingMode
		return nil
	}
}

func WithReclaimPolicy(reclaimPolicy corev1.PersistentVolumeReclaimPolicy) StorageClassOption {
	return func(s *storagev1.StorageClass) error {
		if len(reclaimPolicy) == 0 {
			reclaimPolicy = "Delete"
		}

		s.ReclaimPolicy = &reclaimPolicy
		return nil
	}
}

func WithAllowedTopologies(allowedTopologies map[string][]string) StorageClassOption {
	return func(s *storagev1.StorageClass) error {
		if len(allowedTopologies) == 0 {
			return errors.New("Failed to set AllowedTopologies. " +
				"Input is invalid.")
		}

		appendAllowedTopologies(s, allowedTopologies)
		return nil
	}
}

func WithNodeAffinityLabel(nodeLabelKey string) StorageClassOption {
	return func(s *storagev1.StorageClass) error {
		if len(nodeLabelKey) == 0 {
			return errors.New("Failed to set NodeLabelKey. " +
				"Input is invalid.")
		}

		// Check if the existing parameters and Provisioner name
		// are usable with NodeAffnityLabel.
		// NodeAffinityLabel is only compatible with
		// Hostpath StorageType.
		if !isCompatibleWithNodeAffinityLabel(s) {
			return errors.New("Failed to set NodeAffinityLabel. " +
				"Invalid existing '" + string(mconfig.CASConfigKey) +
				"' annotaion parameters or Provisioner name.")
		}

		config := "- name: NodeAffinityLabel\n" +
			"  value: \"" + nodeLabelKey + "\"\n"

		ok := writeOrAppendCASConfig(s, config)
		if !ok {
			return errors.New("Failed to set NodeAffinityLabel" +
				" parameter for Hostpath.")
		}
		return nil
	}
}

func WithFSType(filesystem string) StorageClassOption {
	return func(s *storagev1.StorageClass) error {
		if !isValidFilesystem(filesystem) {
			return errors.New("Filesystem is invalid. " +
				"Accepted values are \"ext4\" and \"xfs\".")
		}

		// Check if the existing parameters and
		// Provisioner name are usable with FSType.
		// FSType is only compatible with
		// Device StorageType.
		if !isCompatibleWithFSType(s) {
			return errors.New("Failed to set FSType. " +
				"Invalid existing '" + string(mconfig.CASConfigKey) +
				"' annotation parameters or Provisioner name")
		}

		config := "- name: FSType\n" +
			"  value: \"" + filesystem + "\"\n"

		ok := writeOrAppendCASConfig(s, config)
		if !ok {
			return errors.New("Failed to set FSType" +
				" parameter for Device.")
		}
		return nil
	}
}

func WithBlockDeviceTag(bdLabelValue string) StorageClassOption {
	return func(s *storagev1.StorageClass) error {
		if len(bdLabelValue) == 0 {
			return errors.New("Failed to set BlockDeviceTag. " +
				"Input is invalid.")
		}

		// Check if the existing parameters and Provisioner name
		// are usable with BlockDeviceTag.
		// BlockDeviceTag is only compatible with
		// Device StorageType.
		if !isCompatibleWithBlockDeviceTag(s) {
			return errors.New("Failed to set BlockDeviceTag. " +
				"Invalid existing '" + string(mconfig.CASConfigKey) +
				"' annotaion parameters or Provisioner name.")
		}

		config := "- name: BlockDeviceTag\n" +
			"  value: \"" + bdLabelValue + "\"\n"

		ok := writeOrAppendCASConfig(s, config)
		if !ok {
			return errors.New("Failed to set BlockDeviceTag" +
				" parameter for Device.")
		}
		return nil
	}
}