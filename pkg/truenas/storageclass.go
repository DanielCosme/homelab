package truenas

import (
	"danicos.dev/daniel/go-kube/pkg/kube"
	"danicos.dev/daniel/homelab/pkg/root"
	core "k8s.io/api/core/v1"
	storage "k8s.io/api/storage/v1"
)

var NFSStorageClass storage.StorageClass
var iSCSIStorageClass storage.StorageClass

/*
volblocksize	ZVOL block size	512, 1K, 2K, 4K, 8K, 16K, 32K, 64K, 128K
iscsi.blocksize	iSCSI logical block size	512, 1024, 2048, 4096
iscsi.chapUser	CHAP username	string
iscsi.chapSecret	CHAP password (12-16 chars)	string
iscsi.chapPeerUser	Mutual CHAP peer user	string
iscsi.chapPeerSecret	Mutual CHAP peer password	string
iscsi.initiators	Allowed initiator IQNs	comma-separated
iscsi.networks	Allowed network CIDRs	comma-separated
*/

func init() {
	/*
		From:  https://github.com/truenas/truenas-csi
	*/
	NFSStorageClass = storage.StorageClass{
		TypeMeta:    kube.StorageClassMeta,
		ObjectMeta:  kube.ObjectMeta(root.TrueNASSTorageClassNFS, ""),
		Provisioner: root.TrueNASProvisioner,
		Parameters: map[string]string{
			"protocol":    "nfs",
			"compression": "LZ4",
			// Custom NFS mount options passed to clients
			// "nfs.mountOptions": "hard,nfsvers=4.1",
		},
		ReclaimPolicy:        new(core.PersistentVolumeReclaimRetain),
		VolumeBindingMode:    new(storage.VolumeBindingImmediate),
		AllowVolumeExpansion: new(true),
	}
	iSCSIStorageClass = storage.StorageClass{
		TypeMeta:    kube.StorageClassMeta,
		ObjectMeta:  kube.ObjectMeta(root.TrueNASSTorageClass_iSCSI, ""),
		Provisioner: root.TrueNASProvisioner,
		Parameters: map[string]string{
			"protocol":         "iscsi",
			"compression":      "LZ4",
			"volblocksize":     "16K",
			"iscsi.blocksize":  "4096",
			"iscsi.initiators": "ALL",
			"iscsi.networks":   "0.0.0.0/0",
		},
		ReclaimPolicy:        new(core.PersistentVolumeReclaimDelete),
		VolumeBindingMode:    new(storage.VolumeBindingImmediate),
		AllowVolumeExpansion: new(true),
	}
}

/*
NFS Parameters
Parameter					Description						Example

nfs.hosts					Allowed hosts					10.0.0.0/8,192.168.1.0/24
nfs.networks			Allowed networks			10.0.0.0/8
nfs.mountOptions	Client mount options	hard,nfsvers=4.1
nfs.mapAllUser		NFS user mapping 			(default: root)	postgres
nfs.mapAllGroup		NFS group mapping 		(default: wheel)	postgres
*/
