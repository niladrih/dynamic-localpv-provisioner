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
	coreclient "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type PVCClientConfig struct {
	config *rest.Config
}

type PVCClientOption func(*PVCClientConfig) error

func NewPVCClient(opts ...PVCClientOption) (*coreclient.CoreV1Client, error) {
	clientConfig := &PVCClientConfig{}

	var err error
	for _, opt := range opts {
		err = opt(clientConfig)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to build PersistentVolumeClaim client config.")
		}
	}

	var client *coreclient.CoreV1Client
	client, err = coreclient.NewForConfig(clientConfig.config)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Core V1 Client with REST config.")
	}

	return client, nil
}

func WithKubeconfigPath(kubeconfigPath string) PVCClientOption {
	return func(clientConfig *PVCClientConfig) error {
		if len(kubeconfigPath) == 0 {
			return errors.New("Kubeconfig path is empty.")
		}

		var err error
		clientConfig.config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			return errors.Wrap(err, "Failed to get config from kubeconfig path.")
		}

		return nil
	}
}
