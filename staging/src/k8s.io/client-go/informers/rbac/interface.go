/*
Copyright The Kubernetes Authors.

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

// Code generated by informer-gen. DO NOT EDIT.

package rbac

import (
	internalinterfaces "k8s.io/client-go/informers/internalinterfaces"
	v1 "k8s.io/client-go/informers/rbac/v1"
	v1alpha1 "k8s.io/client-go/informers/rbac/v1alpha1"
	v1beta1 "k8s.io/client-go/informers/rbac/v1beta1"
)

// Interface provides access to each of this group's versions.
type Interface interface {
	// V1 provides access to shared informers for resources in V1.
	V1() v1.Interface
	// V1alpha1 provides access to shared informers for resources in V1alpha1.
	V1alpha1() v1alpha1.Interface
	// V1beta1 provides access to shared informers for resources in V1beta1.
	V1beta1() v1beta1.Interface
}

type group struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tenant           string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc, optional_tenant ...string) Interface {
	tenant := "default"
	if len(optional_tenant) > 0 {
		tenant = optional_tenant[0]
	}

	return &group{factory: f, namespace: namespace, tenant: tenant, tweakListOptions: tweakListOptions}
}

// V1 returns a new v1.Interface.
func (g *group) V1() v1.Interface {
	return v1.New(g.factory, g.namespace, g.tweakListOptions, g.tenant)
}

// V1alpha1 returns a new v1alpha1.Interface.
func (g *group) V1alpha1() v1alpha1.Interface {
	return v1alpha1.New(g.factory, g.namespace, g.tweakListOptions, g.tenant)
}

// V1beta1 returns a new v1beta1.Interface.
func (g *group) V1beta1() v1beta1.Interface {
	return v1beta1.New(g.factory, g.namespace, g.tweakListOptions, g.tenant)
}
