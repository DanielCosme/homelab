package truenas

import (
	"fmt"
	"slices"
	"strings"

	"danicos.dev/daniel/go-kube/pkg/kube"
	"danicos.dev/daniel/go-kube/pkg/stack"
	"danicos.dev/daniel/homelab/pkg/root"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	storage "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var Secret = struct {
	Name   string
	APIKey string
}{
	Name:   "truenas-api-credentials",
	APIKey: "api-key",
}

var Config = struct {
	TruenasURL      string
	TrueNASInsecure string
	DefaultPool     string
	NFSServer       string
	ISCSIPortal     string
	ISCSIIQNBase    string
}{
	TruenasURL:      "truenasURL",
	TrueNASInsecure: "truenasInsecure",
	DefaultPool:     "defaultPool",
	NFSServer:       "nfsServer",
	ISCSIPortal:     "iscsiPortal",
	ISCSIIQNBase:    "iscsiIQNBase",
}

var (
	// Origin: https://github.com/truenas/truenas-csi
	Namespace    core.Namespace
	meta         kube.Metadata
	controllerSA core.ServiceAccount
	nodeSA       core.ServiceAccount
	config       core.ConfigMap
)

func init() {
	Namespace = kube.Namespace(root.TrueNAS_CSI)
	meta = kube.NewMetadata(root.TrueNAS_CSI, Namespace)
	controllerSA = kube.ServiceAccount(root.TrueNAS_CSI+"-controller", Namespace)
	nodeSA = kube.ServiceAccount(root.TrueNAS_CSI+"-node", Namespace)
	split := strings.Split(root.TrueNASURL, ".")
	slices.Reverse(split)
	config_meta := kube.NewMetadata("truenas-csi-config", Namespace)
	config = core.ConfigMap{
		TypeMeta:   kube.ConfigMapMeta,
		ObjectMeta: config_meta.Meta(),
		Data: map[string]string{
			Config.TruenasURL:      fmt.Sprintf("wss://%s/api/current", root.TRUENAS_IP),
			Config.TrueNASInsecure: "true",
			Config.DefaultPool:     "datapool",
			Config.NFSServer:       root.TRUENAS_IP,
			Config.ISCSIPortal:     fmt.Sprintf("%s:3260", root.TRUENAS_IP),
			// Config.ISCSIIQNBase:    fmt.Sprintf("iqn.2026-04.%s", strings.Join(split, ".")),
			// apex-truenas.orca-uaru.ts.net
			// 10.0.0.57
		},
	}
}

func Stack() stack.Stack {
	kz := kube.NewKuztomizedStack(
		meta,
		map[string]any{
			"namespace": Namespace,
			// "controller-deployment":      controllerDeployment(),
			// "controller-service-account": controllerSA,
			// "controller-cluster-role":    controllerRole,
			// "controller-binding":         kube.ClusterRoleBinding(controllerRole.Name+"-binding", controllerSA, controllerRole),
			// "node-service-account":       nodeSA,
			// "node-cluster-role":          nodeRole,
			// "node-binding":               kube.ClusterRoleBinding(nodeRole.Name+"-binding", nodeSA, nodeRole),
			// "node-deamonset": nodeCSI(),
			// "CSIDriver":                  CSIDriver(root.TrueNASProvisioner),
			"config":            config,
			"nfs-storage-class": NFSStorageClass,
			// "iscsi-storage-class": iSCSIStorageClass,
		},
	)
	return kz.Stack(root.TrueNAS_CSI)
}

func controllerDeployment() apps.Deployment {
	meta := kube.NewMetadata(root.TrueNAS_CSI+"-controller", Namespace)
	vol := core.Volume{
		Name: "socket-dir",
		VolumeSource: core.VolumeSource{
			EmptyDir: &core.EmptyDirVolumeSource{},
		},
	}
	spec := core.PodSpec{
		ServiceAccountName: controllerSA.Name,
		Containers: []core.Container{
			{
				Name:            "csi-controller",
				Image:           "ghcr.io/truenas/truenas-csi:latest",
				ImagePullPolicy: core.PullIfNotPresent,
				Args: []string{
					"--endpoint=$(CSI_ENDPOINT)",
					"--node-id=$(NODE_ID)",
					"--mode=controller",
					"--v=4",
				},
				Env:          controllerEnv(),
				VolumeMounts: []core.VolumeMount{{Name: vol.Name, MountPath: "/csi"}},
				LivenessProbe: &core.Probe{
					ProbeHandler: core.ProbeHandler{
						HTTPGet: &core.HTTPGetAction{
							Path: "/healthz",
							Port: intstr.FromInt(9808),
						},
					},
					InitialDelaySeconds: 10,
					TimeoutSeconds:      3,
					PeriodSeconds:       10,
					FailureThreshold:    5,
				},
			},
			{
				Name:  "csi-provisioner",
				Image: "registry.k8s.io/sig-storage/csi-provisioner:v6.1.1",
				Args: []string{
					"--csi-address=$(ADDRESS)",
					"--v=5",
					"--feature-gates=Topology=true",
					"--extra-create-metadata",
					"--leader-election=true",
					"--default-fstype=ext4",
					"--timeout=60s",
				},
				Env: []core.EnvVar{{
					Name:  "ADDRESS",
					Value: "/csi/csi.sock",
				}},
				VolumeMounts: []core.VolumeMount{{Name: vol.Name, MountPath: "/csi"}},
			},
			{
				Name:  "csi-attacher",
				Image: "registry.k8s.io/sig-storage/csi-attacher:v4.11.0",
				Args: []string{
					"--csi-address=$(ADDRESS)",
					"--v=5",
					"--leader-election=true",
					"--timeout=60s",
				},
				Env: []core.EnvVar{{
					Name:  "ADDRESS",
					Value: "/csi/csi.sock",
				}},
				VolumeMounts: []core.VolumeMount{{Name: vol.Name, MountPath: "/csi"}},
			},
			{
				Name:  "csi-snapshotter",
				Image: "registry.k8s.io/sig-storage/csi-snapshotter:v8.5.0",
				Args: []string{
					"--csi-address=$(ADDRESS)",
					"--v=5",
					"--leader-election=true",
					"--timeout=60s",
				},
				Env: []core.EnvVar{{
					Name:  "ADDRESS",
					Value: "/csi/csi.sock",
				}},
				VolumeMounts: []core.VolumeMount{{Name: vol.Name, MountPath: "/csi"}},
			},
			{
				Name:  "csi-resizer",
				Image: "registry.k8s.io/sig-storage/csi-resizer:v2.1.0",
				Args: []string{
					"--csi-address=$(ADDRESS)",
					"--v=5",
					"--leader-election=true",
					"--timeout=60s",
				},
				Env: []core.EnvVar{{
					Name:  "ADDRESS",
					Value: "/csi/csi.sock",
				}},
				VolumeMounts: []core.VolumeMount{{Name: vol.Name, MountPath: "/csi"}},
			},
			{
				Name:  "liveness-probe",
				Image: "registry.k8s.io/sig-storage/livenessprobe:v2.18.0",
				Args: []string{
					"--csi-address=/csi/csi.sock",
					"--health-port=9808",
				},
				VolumeMounts: []core.VolumeMount{{Name: vol.Name, MountPath: "/csi"}},
			},
		},
		Volumes: []core.Volume{vol},
	}
	return kube.NewDeployment(meta, spec)
}

func controllerEnv() []core.EnvVar {
	envMapping := map[string]string{
		"CSI_ENDPOINT": "unix:///csi/csi.sock",
	}
	secretMapping := map[string]string{
		"TRUENAS_API_KEY": Secret.APIKey,
	}
	env1 := kube.NewEnvVarWithSecret(envMapping, secretMapping, Secret.Name)
	envMapping = map[string]string{
		"TRUENAS_URL":                  Config.TruenasURL,
		"TRUENAS_DEFAULT_POOL":         Config.DefaultPool,
		"TRUENAS_NFS_SERVER":           Config.NFSServer,
		"TRUENAS_ISCSI_PORTAL":         Config.ISCSIPortal,
		"TRUENAS_ISCSI_IQN_BASE":       Config.ISCSIIQNBase,
		"TRUENAS_INSECURE_SKIP_VERIFY": Config.TrueNASInsecure,
	}
	env2 := kube.NewEnvVarWithConfig(envMapping, config)
	nodeEnv := core.EnvVar{
		Name: "NODE_ID", // value
		ValueFrom: &core.EnvVarSource{
			FieldRef: &core.ObjectFieldSelector{
				FieldPath: "spec.nodeName",
			},
		},
	}
	env3 := append(env1, nodeEnv)
	return append(env3, env2...)
}

func CSIDriver(name string) storage.CSIDriver {
	fsGroupPolicy := storage.FileFSGroupPolicy
	spec := storage.CSIDriverSpec{
		AttachRequired: new(true),
		PodInfoOnMount: new(true),
		VolumeLifecycleModes: []storage.VolumeLifecycleMode{
			storage.VolumeLifecycleEphemeral,
			storage.VolumeLifecyclePersistent,
		},
		FSGroupPolicy: &fsGroupPolicy,
	}
	return kube.CSIDriver(name, spec)
}

func nodeCSI() apps.DaemonSet {
	registrationDir := core.Volume{
		Name: "registration-dir",
		VolumeSource: core.VolumeSource{
			HostPath: &core.HostPathVolumeSource{
				Path: "/var/lib/kubelet/plugins_registry/",
				Type: new(core.HostPathDirectoryOrCreate),
			},
		},
	}
	pluginDir := core.Volume{
		Name: "plugin-dir",
		VolumeSource: core.VolumeSource{
			HostPath: &core.HostPathVolumeSource{
				Path: "/var/lib/kubelet/plugins/csi.truenas.io/",
				Type: new(core.HostPathDirectoryOrCreate),
			},
		},
	}
	kubeletDir := core.Volume{
		Name: "kubelet-dir",
		VolumeSource: core.VolumeSource{
			HostPath: &core.HostPathVolumeSource{
				Path: "/var/lib/kubelet",
				Type: new(core.HostPathDirectory),
			},
		},
	}
	deviceDir := core.Volume{
		Name: "device-dir",
		VolumeSource: core.VolumeSource{
			HostPath: &core.HostPathVolumeSource{
				Path: "/dev",
			},
		},
	}
	iscsiDir := core.Volume{
		Name: "iscsi-dir",
		VolumeSource: core.VolumeSource{
			HostPath: &core.HostPathVolumeSource{
				Path: "/etc/iscsi",
				Type: new(core.HostPathDirectory),
			},
		},
	}
	iscsiLib := core.Volume{
		Name: "iscsi-lib",
		VolumeSource: core.VolumeSource{
			HostPath: &core.HostPathVolumeSource{
				Path: "/var/lib/iscsi",
				Type: new(core.HostPathDirectoryOrCreate),
			},
		},
	}
	hostRoot := core.Volume{
		Name: "host-root",
		VolumeSource: core.VolumeSource{
			HostPath: &core.HostPathVolumeSource{
				Path: "/",
				Type: new(core.HostPathDirectory),
			},
		},
	}
	meta := kube.NewMetadata(root.TrueNAS_CSI+"-node", Namespace)
	podSpec := core.PodSpec{
		ServiceAccountName: nodeSA.Name,
		HostNetwork:        true,
		HostPID:            true,
		HostIPC:            true,
		DNSPolicy:          core.DNSClusterFirstWithHostNet,
		PriorityClassName:  "system-node-critical",
		Tolerations: []core.Toleration{{
			Operator: core.TolerationOpExists,
		}},
		Containers: []core.Container{
			{
				Name:            "csi-node",
				Image:           "ghcr.io/truenas/truenas-csi:latest",
				ImagePullPolicy: core.PullIfNotPresent,
				Lifecycle: &core.Lifecycle{
					PostStart: &core.LifecycleHandler{
						Exec: &core.ExecAction{
							Command: []string{
								"/bin/sh",
								"-c",
								"mkdir -p /run/lock/iscsi && mv /usr/sbin/iscsiadm /usr/sbin/iscsiadm.orig 2>/dev/null; printf '#!/bin/sh\\nnsenter --mount=/host/proc/1/ns/mnt -- /usr/sbin/iscsiadm \"$@\"\\n' > /usr/sbin/iscsiadm && chmod +x /usr/sbin/iscsiadm",
							},
						},
					},
				},
				SecurityContext: &core.SecurityContext{Privileged: new(true)},
				Args: []string{
					"--endpoint=$(CSI_ENDPOINT)",
					"--node-id=$(NODE_ID)",
					"--mode=node",
					"--v=4",
				},
				Env: controllerEnv(),
				VolumeMounts: []core.VolumeMount{
					{
						Name:      pluginDir.Name,
						MountPath: "/csi",
					},
					{
						Name:      kubeletDir.Name,
						MountPath: "/var/lib/kubelet",
					},
					{
						Name:      deviceDir.Name,
						MountPath: "/dev",
					},
					{
						Name:             iscsiDir.Name,
						MountPath:        "/etc/iscsi",
						MountPropagation: new(core.MountPropagationBidirectional),
					},
					{
						Name:             iscsiLib.Name,
						MountPath:        "/var/lib/iscsi",
						MountPropagation: new(core.MountPropagationBidirectional),
					},
					{
						Name:             hostRoot.Name,
						MountPath:        "/host",
						MountPropagation: new(core.MountPropagationBidirectional),
					},
				},
				LivenessProbe: &core.Probe{
					ProbeHandler: core.ProbeHandler{
						HTTPGet: &core.HTTPGetAction{
							Path: "/healthz",
							Port: intstr.FromInt(9808),
						},
					},
					InitialDelaySeconds: 10,
					TimeoutSeconds:      3,
					PeriodSeconds:       10,
					FailureThreshold:    5,
				},
			},
			{
				Name:  "csi-node-driver-registrar",
				Image: "registry.k8s.io/sig-storage/csi-node-driver-registrar:v2.16.0",
				Args: []string{
					"--csi-address=$(ADDRESS)",
					"--kubelet-registration-path=$(DRIVER_REG_SOCK_PATH)",
					"--v=5",
				},
				Env: []core.EnvVar{
					{
						Name:  "ADDRESS",
						Value: "/csi/csi.sock",
					},
					{
						Name:  "DRIVER_REG_SOCK_PATH",
						Value: "/var/lib/kubelet/plugins/csi.truenas.io/csi.sock",
					},
				},
				VolumeMounts: []core.VolumeMount{
					{
						Name:      pluginDir.Name,
						MountPath: "/csi",
					},
					{
						Name:      registrationDir.Name,
						MountPath: "/registration",
					},
				},
			},
			{
				Name:  "liveness-probe",
				Image: "registry.k8s.io/sig-storage/livenessprobe:v2.18.0",
				Args: []string{
					"--csi-address=/csi/csi.sock",
					"--health-port=9808",
				},
				VolumeMounts: []core.VolumeMount{{Name: pluginDir.Name, MountPath: "/csi"}},
			},
		},
		Volumes: []core.Volume{
			registrationDir,
			pluginDir,
			kubeletDir,
			deviceDir,
			iscsiDir,
			iscsiLib,
			hostRoot,
		},
	}

	return kube.NewDeamonSet(meta, podSpec)
}
