/*
Copyright 2017 The Gardener Authors.

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
package driver

import (
	"github.com/gardener/node-controller-manager/pkg/apis/machine/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

//Driver is the common interface for creation/deletion of the VMs over different cloud-providers.
type Driver interface {
	Create() (string, string, error)
	Delete() error
	GetExisting() (string, error)
}

func NewDriver(machineId string, secretRef *corev1.Secret, classKind string, machineClass interface{}, machineName string) Driver {

	switch classKind {
	case "AWSMachineClass":
		return &AWSDriver{
			AWSMachineClass: machineClass.(*v1alpha1.AWSMachineClass),
			CloudConfig:     secretRef,
			UserData:        string(secretRef.Data["userData"]),
			MachineId:       machineId,
			MachineName:     machineName,
		}

	case "AzureMachineClass":
		return &AzureDriver{
			AzureMachineClass: machineClass.(*v1alpha1.AzureMachineClass),
			CloudConfig:       secretRef,
			UserData:          string(secretRef.Data["userData"]),
			MachineId:         machineId,
			MachineName:       machineName,
		}

	case "GCPMachineClass":
		return &GCPDriver{
			GCPMachineClass: machineClass.(*v1alpha1.GCPMachineClass),
			CloudConfig:     secretRef,
			UserData:        string(secretRef.Data["userData"]),
			MachineId:       machineId,
			MachineName:     machineName,
		}
	}

	return NewFakeDriver(
		func() (string, string, error) {
			return "fake", "fake_ip", nil
		},
		func() error {
			return nil
		},
		func() (string, error) {
			return "fake", nil
		},
	)
}
