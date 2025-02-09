/*
Copyright 2021 The Kubernetes Authors.

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

package nodegroupset

import (
	"k8s.io/autoscaler/cluster-autoscaler/config"
	schedulerframework "k8s.io/kubernetes/pkg/scheduler/framework"
)

// CreateClusterAPINodeInfoComparator returns a comparator that checks if two nodes should be considered
// part of the same NodeGroupSet. This is true if they match usual conditions checked by IsCloudProviderNodeInfoSimilar,
// even if they have different infrastructure provider-specific labels.
func CreateClusterAPINodeInfoComparator(extraIgnoredLabels []string, ratioOpts config.NodeGroupDifferenceRatios) NodeInfoComparator {
	capiIgnoredLabels := map[string]bool{
		"topology.ebs.csi.aws.com/zone":                 true, // this is a label used by the AWS EBS CSI driver as a target for Persistent Volume Node Affinity
		"topology.diskplugin.csi.alibabacloud.com/zone": true, // this is a label used by the Alibaba Cloud CSI driver as a target for Persistent Volume Node Affinity
		"ibm-cloud.kubernetes.io/worker-id":             true, // this is a label used by the IBM Cloud Cloud Controler Manager
		"vpc-block-csi-driver-labels":                   true, // this is a label used by the IBM Cloud CSI driver as a target for Persisten Volume Node Affinity

	}

	for k, v := range BasicIgnoredLabels {
		capiIgnoredLabels[k] = v
	}

	for _, k := range extraIgnoredLabels {
		capiIgnoredLabels[k] = true
	}

	return func(n1, n2 *schedulerframework.NodeInfo) bool {
		return IsCloudProviderNodeInfoSimilar(n1, n2, capiIgnoredLabels, ratioOpts)
	}
}
