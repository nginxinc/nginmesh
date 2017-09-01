// The code in this file reuses the code from https://github.com/istio/pilot/blob/pilot-0-2-0-working/platform/kube/inject/inject.go
// licensed under the following conditions:
// Copyright 2017 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"fmt"
	"io"

	"github.com/ghodss/yaml"
	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	yamlDecoder "k8s.io/apimachinery/pkg/util/yaml"
)

type params struct {
	InitImage       string
	ProxyImage      string
	SidecarProxyUID int64
	ProxyListenPort int
}

func injectIntoPodTemplateSpec(p *params, t *v1.PodTemplateSpec) error {
	initContainer := v1.Container{
		Name:            "init",
		Image:           p.InitImage,
		ImagePullPolicy: v1.PullAlways,
		SecurityContext: &v1.SecurityContext{
			Capabilities: &v1.Capabilities{
				Add: []v1.Capability{"NET_ADMIN"},
			},
		},
		Args: []string{"-p", fmt.Sprint(p.ProxyListenPort), "-u", fmt.Sprint(p.SidecarProxyUID)},
	}

	proxyContainer := v1.Container{
		Name:            "proxy",
		Image:           p.ProxyImage,
		ImagePullPolicy: v1.PullAlways,
		SecurityContext: &v1.SecurityContext{
			RunAsUser: &p.SidecarProxyUID,
		},
		Env: []v1.EnvVar{{
			Name: "POD_NAME",
			ValueFrom: &v1.EnvVarSource{
				FieldRef: &v1.ObjectFieldSelector{
					FieldPath: "metadata.name",
				},
			},
		}, {
			Name: "POD_NAMESPACE",
			ValueFrom: &v1.EnvVarSource{
				FieldRef: &v1.ObjectFieldSelector{
					FieldPath: "metadata.namespace",
				},
			},
		}, {
			Name: "POD_IP",
			ValueFrom: &v1.EnvVarSource{
				FieldRef: &v1.ObjectFieldSelector{
					FieldPath: "status.podIP",
				},
			},
		}},
	}

	t.Spec.InitContainers = append(t.Spec.InitContainers, initContainer)
	t.Spec.Containers = append(t.Spec.Containers, proxyContainer)
	return nil
}

func injectIntoResourceFile(p *params, in io.Reader, out io.Writer) error {
	reader := yamlDecoder.NewYAMLReader(bufio.NewReaderSize(in, 4096))
	for {
		raw, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		kinds := map[string]struct {
			typ    interface{}
			inject func(typ interface{}) error
		}{
			"DaemonSet": {
				typ: &v1beta1.DaemonSet{},
				inject: func(typ interface{}) error {
					return injectIntoPodTemplateSpec(p, &((typ.(*v1beta1.DaemonSet)).Spec.Template))
				},
			},
			"ReplicaSet": {
				typ: &v1beta1.ReplicaSet{},
				inject: func(typ interface{}) error {
					return injectIntoPodTemplateSpec(p, &((typ.(*v1beta1.ReplicaSet)).Spec.Template))
				},
			},
			"Deployment": {
				typ: &v1beta1.Deployment{},
				inject: func(typ interface{}) error {
					return injectIntoPodTemplateSpec(p, &((typ.(*v1beta1.Deployment)).Spec.Template))
				},
			},
			"ReplicationController": {
				typ: &v1.ReplicationController{},
				inject: func(typ interface{}) error {
					return injectIntoPodTemplateSpec(p, ((typ.(*v1.ReplicationController)).Spec.Template))
				},
			},
		}
		var updated []byte
		var meta metav1.TypeMeta
		if err = yaml.Unmarshal(raw, &meta); err != nil {
			return err
		}
		if kind, ok := kinds[meta.Kind]; ok {
			if err = yaml.Unmarshal(raw, kind.typ); err != nil {
				return err
			}
			if err = kind.inject(kind.typ); err != nil {
				return err
			}
			if updated, err = yaml.Marshal(kind.typ); err != nil {
				return err
			}
		} else {
			updated = raw // unchanged
		}

		if _, err = out.Write(updated); err != nil {
			return err
		}
		if _, err = fmt.Fprint(out, "---\n"); err != nil {
			return err
		}
	}
	return nil
}
