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

package v1

import (
	time "time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	internalinterfaces "k8s.io/client-go/informers/internalinterfaces"
	kubernetes "k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/listers/core/v1"
	cache "k8s.io/client-go/tools/cache"
)

// ReplicationControllerInformer provides access to a shared informer and lister for
// ReplicationControllers.
type ReplicationControllerInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.ReplicationControllerLister
}

type replicationControllerInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
	tenant           string
}

// NewReplicationControllerInformer constructs a new informer for ReplicationController type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewReplicationControllerInformer(client kubernetes.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, optional_tenant ...string) cache.SharedIndexInformer {
	return NewFilteredReplicationControllerInformer(client, namespace, resyncPeriod, indexers, nil, optional_tenant...)
}

// NewFilteredReplicationControllerInformer constructs a new informer for ReplicationController type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredReplicationControllerInformer(client kubernetes.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc, optional_tenant ...string) cache.SharedIndexInformer {
	tenant := "default"
	if len(optional_tenant) > 0 {
		tenant = optional_tenant[0]
	}

	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CoreV1().ReplicationControllers(namespace, tenant).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CoreV1().ReplicationControllers(namespace, tenant).Watch(options)
			},
		},
		&corev1.ReplicationController{},
		resyncPeriod,
		indexers,
	)
}

func (f *replicationControllerInformer) defaultInformer(client kubernetes.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredReplicationControllerInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions, f.tenant)
}

func (f *replicationControllerInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&corev1.ReplicationController{}, f.defaultInformer)
}

func (f *replicationControllerInformer) Lister() v1.ReplicationControllerLister {
	return v1.NewReplicationControllerLister(f.Informer().GetIndexer())
}
