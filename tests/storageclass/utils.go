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
	"path/filepath"
	"strings"

	cast "github.com/openebs/maya/pkg/castemplate/v1alpha1"
)

func isValidPath(hostpath string) bool {
	// Is an abolute path
	if !filepath.IsAbs(hostpath) {
		return false
	}

	// IsNotRoot
	path := strings.TrimSuffix(string(hostpath), "/")
	parentDir, subDir := filepath.Split(path)
	parentDir = strings.TrimSuffix(parentDir, "/")
	subDir = strings.TrimSuffix(subDir, "/")
	if parentDir == "" || subDir == "" {
		return false
	}
	return true
}

// Used to check if the cas.openebs.io/config value string
// has valid parameters for hostpath or not
// e.g.
// Parameters like 'BlockDeviceTag', already existing 'StorageType'
// are incompatible.
func isCompatibleWithHostpath(scCASConfigStr string) bool {
	// Unmarshall to mconfig.Config
	scCASConfig, err := cast.UnMarshallToConfig(scCASConfigStr)
	if err != nil {
		return false
	}

	// Check for invalid CAS config parameters
	for _, config := range scCASConfig {
		switch strings.TrimSpace(config.Name) {
		case "NodeAffinityLabel":
			continue
		default:
			return false
		}
	}
	return true
}

func isCompatibleWithNodeAffinityLabel(scCASConfigStr string) bool {
	// Unmarshall to mconfig.Config
	scCASConfig, err := cast.UnMarshallToConfig(scCASConfigStr)
	if err != nil {
		return false
	}

	// Check for invalid CAS config parameters
	for _, config := range scCASConfig {
		switch strings.TrimSpace(config.Name) {
		case "StorageType":
			if config.Value == "\"hostpath\"" || config.Value == "hostpath" {
				continue
			} else {
				return false
			}
		case "BasePath":
			if !isValidPath(config.Value) {
				return false
			}
			continue
		default:
			return false
		}
	}
	return true
}

func isCompatibleWithDevice(scCASConfigStr string) bool {
	// Unmarshall to mconfig.Config
	scCASConfig, err := cast.UnMarshallToConfig(scCASConfigStr)
	if err != nil {
		return false
	}

	// Check for invalid CAS config parameters
	for _, config := range scCASConfig {
		switch strings.TrimSpace(config.Name) {
		case "BlockDeviceTag":
			continue
		case "FSType":
			continue
		default:
			return false
		}
	}
	return true
}

func isCompatibleWithFSType(scCASConfigStr string) bool {
	// Unmarshall to mconfig.Config
	scCASConfig, err := cast.UnMarshallToConfig(scCASConfigStr)
	if err != nil {
		return false
	}

	// Check for invalid CAS config parameters
	for _, config := range scCASConfig {
		switch strings.TrimSpace(config.Name) {
		case "StorageType":
			if config.Value == "\"device\"" || config.Value == "device" {
				continue
			} else {
				return false
			}
		case "BlockDeviceTag":
			continue
		default:
			return false
		}
	}
	return true
}

func isValidFilesystem(filesystem string) bool {
	switch filesystem {
	case "xfs":
		return true
	case "ext4":
		return true
	default:
		return false
	}
}

func isCompatibleWithBlockDeviceTag(scCASConfigStr string) bool {
	// Unmarshall to mconfig.Config
	scCASConfig, err := cast.UnMarshallToConfig(scCASConfigStr)
	if err != nil {
		return false
	}

	// Check for invalid CAS config parameters
	for _, config := range scCASConfig {
		switch strings.TrimSpace(config.Name) {
		case "StorageType":
			if config.Value == "\"device\"" || config.Value == "device" {
				continue
			} else {
				return false
			}
		case "FSType":
			if !isValidFilesystem(config.Value) {
				return false
			}
			continue
		default:
			return false
		}
	}
	return true
}
