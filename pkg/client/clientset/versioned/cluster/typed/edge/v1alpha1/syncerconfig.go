//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright The KubeStellar Authors.

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

// Code generated by kcp code-generator. DO NOT EDIT.

package v1alpha1

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	kcpclient "github.com/kcp-dev/apimachinery/v2/pkg/client"
	"github.com/kcp-dev/logicalcluster/v3"

	edgev1alpha1 "github.com/kubestellar/kubestellar/pkg/apis/edge/v1alpha1"
	edgev1alpha1client "github.com/kubestellar/kubestellar/pkg/client/clientset/versioned/typed/edge/v1alpha1"
)

// SyncerConfigsClusterGetter has a method to return a SyncerConfigClusterInterface.
// A group's cluster client should implement this interface.
type SyncerConfigsClusterGetter interface {
	SyncerConfigs() SyncerConfigClusterInterface
}

// SyncerConfigClusterInterface can operate on SyncerConfigs across all clusters,
// or scope down to one cluster and return a edgev1alpha1client.SyncerConfigInterface.
type SyncerConfigClusterInterface interface {
	Cluster(logicalcluster.Path) edgev1alpha1client.SyncerConfigInterface
	List(ctx context.Context, opts metav1.ListOptions) (*edgev1alpha1.SyncerConfigList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
}

type syncerConfigsClusterInterface struct {
	clientCache kcpclient.Cache[*edgev1alpha1client.EdgeV1alpha1Client]
}

// Cluster scopes the client down to a particular cluster.
func (c *syncerConfigsClusterInterface) Cluster(clusterPath logicalcluster.Path) edgev1alpha1client.SyncerConfigInterface {
	if clusterPath == logicalcluster.Wildcard {
		panic("A specific cluster must be provided when scoping, not the wildcard.")
	}

	return c.clientCache.ClusterOrDie(clusterPath).SyncerConfigs()
}

// List returns the entire collection of all SyncerConfigs across all clusters.
func (c *syncerConfigsClusterInterface) List(ctx context.Context, opts metav1.ListOptions) (*edgev1alpha1.SyncerConfigList, error) {
	return c.clientCache.ClusterOrDie(logicalcluster.Wildcard).SyncerConfigs().List(ctx, opts)
}

// Watch begins to watch all SyncerConfigs across all clusters.
func (c *syncerConfigsClusterInterface) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientCache.ClusterOrDie(logicalcluster.Wildcard).SyncerConfigs().Watch(ctx, opts)
}
