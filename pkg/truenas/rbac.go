package truenas

import (
	"danicos.dev/daniel/go-kube/pkg/kube"
	"danicos.dev/daniel/homelab/pkg/root"

	rbac "k8s.io/api/rbac/v1"
)

func controllerClusterRole() rbac.ClusterRole {
	verbsReadUpdate := append(kube.VerbsRead(), kube.VerbsMutate()...)
	rules := []rbac.PolicyRule{
		kube.PolicyRule(kube.APIGroupCore, kube.ResourcePVs, kube.VerbsAll()),
		kube.PolicyRule(kube.APIGroupCore, kube.ResourcePVCs, verbsReadUpdate),
		kube.PolicyRule(kube.APIGroupCore, kube.ResourcePVCsStatus, kube.VerbsMutate()),
		kube.PolicyRule(kube.APIGroupCore, kube.ResourceEvents, verbsReadUpdate),
		kube.PolicyRule(kube.APIGroupCore, kube.ResourceNodes, kube.VerbsRead()),
		kube.PolicyRule(kube.APIGroupCore, kube.ResourcePods, kube.VerbsRead()),
		kube.PolicyRule(kube.APIGroupStorage, kube.ResourceStorageClasses, kube.VerbsRead()),
		kube.PolicyRule(kube.APIGroupStorage, kube.ResourceCSINodes, kube.VerbsRead()),
		kube.PolicyRule(kube.APIGroupStorage, kube.ResourceVolumeAttachments, kube.VerbsAll()),
		kube.PolicyRule(kube.APIGroupStorage, kube.ResourceVolumeAttachmentsStatus, []string{kube.VerbPatch}),
		kube.PolicyRule(kube.APIGroupSnapshot, kube.ResourceVolumeSnapshots, verbsReadUpdate),
		kube.PolicyRule(kube.APIGroupSnapshot, kube.ResourceVolumeSnapshotsStatus, kube.VerbsMutate()),
		kube.PolicyRule(kube.APIGroupSnapshot, kube.ResourceVolumeSnapshotContents, kube.VerbsAll()),
		kube.PolicyRule(kube.APIGroupSnapshot, kube.ResourceVolumeSnapshotContentsStatus, kube.VerbsMutate()),
		kube.PolicyRule(kube.APIGroupSnapshot, kube.ResourceVolumeSnapshotClases, kube.VerbsRead()),
	}
	return kube.ClusterRole(root.TrueNAS_CSI+"-controller-role", rules)
}

func nodeClusterRole() rbac.ClusterRole {
	rules := []rbac.PolicyRule{
		kube.PolicyRule(kube.APIGroupCore, kube.ResourceNodes, []string{kube.VerbGet}),
		kube.PolicyRule(kube.APIGroupCore, kube.ResourcePods, kube.VerbsRead()),
		kube.PolicyRule(kube.APIGroupStorage, kube.ResourceVolumeAttachments, kube.VerbsRead()),
	}
	return kube.ClusterRole(root.TrueNAS_CSI+"-node-role", rules)
}
