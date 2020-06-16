package operator

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/blang/semver"
	csvv1alpha1 "github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"
	csvVersion "github.com/operator-framework/operator-lifecycle-manager/pkg/lib/version"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const openstackClusterName = "openstack-cluster-operator"

func getDeploymentSpec(namespace, image, imagePullPolicy string) appsv1.DeploymentSpec {
	return appsv1.DeploymentSpec{
		Replicas: int32Ptr(1),
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"name": openstackClusterName,
			},
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"name": openstackClusterName,
				},
			},
			Spec: corev1.PodSpec{
				ServiceAccountName: openstackClusterName,
				Containers: []corev1.Container{
					{
						Name:            openstackClusterName,
						Image:           image,
						ImagePullPolicy: corev1.PullPolicy(imagePullPolicy),
						Command:         []string{openstackClusterName},
						//ReadinessProbe: &corev1.Probe{
						//	Handler: corev1.Handler{
						//		Exec: &corev1.ExecAction{
						//			Command: []string{
						//				"stat",
						//				"/tmp/operator-sdk-ready",
						//			},
						//		},
						//	},
						//	InitialDelaySeconds: 5,
						//	PeriodSeconds:       5,
						//	FailureThreshold:    1,
						//},
						Env: []corev1.EnvVar{
							{
								Name:  "OPERATOR_IMAGE",
								Value: image,
							},
							{
								Name:  "OPERATOR_NAME",
								Value: openstackClusterName,
							},
							{
								Name:  "OPERATOR_NAMESPACE",
								Value: namespace,
							},
							{
								Name: "POD_NAME",
								ValueFrom: &corev1.EnvVarSource{
									FieldRef: &corev1.ObjectFieldSelector{
										FieldPath: "metadata.name",
									},
								},
							},
							{
								Name:  "WATCH_NAMESPACE",
								Value: "",
							},
						},
					},
				},
			},
		},
	}
}

// GetClusterPermissions returns the cluster policy rules
func GetClusterPermissions() []rbacv1.PolicyRule {
	return []rbacv1.PolicyRule{
		{
			APIGroups: []string{
				"openstackcluster.openstack.org",
			},
			Resources: []string{
				"*",
			},
			Verbs: []string{
				"*",
			},
		},
		{
			APIGroups: []string{
				"",
			},
			Resources: []string{
				"pods",
				"services",
				"services/finalizers",
				"endpoints",
				"persistentvolumeclaims",
				"events",
				"configmaps",
				"secrets",
				"serviceaccounts",
			},
			Verbs: []string{
				"*",
			},
		},
		{
			APIGroups: []string{
				"apps",
			},
			Resources: []string{
				"deployments",
				"deployments/finalizers",
				"daemonsets",
				"replicasets",
			},
			Verbs: []string{
				"get",
				"list",
				"watch",
				"create",
				"delete",
				"update",
			},
		},
		{
			APIGroups: []string{
				"batch",
			},
			Resources: []string{
				"jobs",
			},
			Verbs: []string{
				"get",
				"list",
				"watch",
				"create",
				"delete",
			},
		},
		{
			APIGroups: []string{
				"rbac.authorization.k8s.io",
			},
			Resources: []string{
				"clusterroles",
				"clusterrolebindings",
				"roles",
				"rolebindings",
			},
			Verbs: []string{
				"get",
				"list",
				"watch",
				"create",
				"delete",
			},
		},
		{
			APIGroups: []string{
				"apiextensions.k8s.io",
			},
			Resources: []string{
				"customresourcedefinitions",
			},
			Verbs: []string{
				"get",
				"list",
				"watch",
				"create",
				"delete",
				"patch",
				"update",
			},
		},
		{
			APIGroups: []string{
				"security.openshift.io",
			},
			Resources: []string{
				"securitycontextconstraints",
			},
			Verbs: []string{
				"get",
				"list",
				"watch",
			},
		},
		{
			APIGroups: []string{
				"security.openshift.io",
			},
			Resources: []string{
				"securitycontextconstraints",
			},
			ResourceNames: []string{
				"privileged",
			},
			Verbs: []string{
				"get",
				"patch",
				"update",
			},
		},
		{
			APIGroups: []string{
				"monitoring.coreos.com",
			},
			Resources: []string{
				"servicemonitors",
			},
			Verbs: []string{
				"get",
				"create",
			},
		},
		{
			APIGroups: []string{
				"operators.coreos.com",
			},
			Resources: []string{
				"clusterserviceversions",
			},
			Verbs: []string{
				"get",
				"list",
				"watch",
			},
		},
	}
}

// GetInstallStrategyBase returns the cluster base strategy
func GetInstallStrategyBase(namespace, image, imagePullPolicy string) csvv1alpha1.StrategyDetailsDeployment {
	return csvv1alpha1.StrategyDetailsDeployment{
		DeploymentSpecs: []csvv1alpha1.StrategyDeploymentSpec{
			csvv1alpha1.StrategyDeploymentSpec{
				Name: "openstack-cluster-operator",
				Spec: getDeploymentSpec(namespace, image, imagePullPolicy),
			},
		},
		Permissions: []csvv1alpha1.StrategyDeploymentPermissions{},
		ClusterPermissions: []csvv1alpha1.StrategyDeploymentPermissions{
			csvv1alpha1.StrategyDeploymentPermissions{
				ServiceAccountName: openstackClusterName,
				Rules:              GetClusterPermissions(),
			},
		},
	}
}

// GetCSVBase returns a base OpenStack Cluster CSV without an InstallStrategy
func GetCSVBase(name, namespace, displayName, description, image, replaces string, version semver.Version, crdDisplay string) *csvv1alpha1.ClusterServiceVersion {
	almExamples, _ := json.Marshal([]interface{}{
		map[string]interface{}{
			"apiVersion": "openstackcluster.openstack.org/v1",
			"kind":       "OpenStackCluster",
			"metadata": map[string]string{
				"name":      "openstack-cluster-operator",
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"BareMetalPlatform": false,
			},
		},
	})

	return &csvv1alpha1.ClusterServiceVersion{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "operators.coreos.com/v1alpha1",
			Kind:       "ClusterServiceVersion",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%v.v%v", name, version.String()),
			Namespace: "placeholder",
			Annotations: map[string]string{
				"alm-examples":   string(almExamples),
				"capabilities":   "Basic Install",
				"certified":      "false",
				"categories":     "OpenShift Optional",
				"containerImage": image,
				"createdAt":      time.Now().Format("2006-01-02 15:04:05"),
				"description":    description,
				"repository":     "https://github.com/openstack-k8s-operators/openstack-cluster-operator",
				"support":        "false",
				"operatorframework.io/suggested-namespace": namespace,
			},
		},
		Spec: csvv1alpha1.ClusterServiceVersionSpec{
			DisplayName: displayName,
			Description: description,
			Keywords:    []string{"OpenStack"},
			Version: csvVersion.OperatorVersion{
				Version: version},
			Replaces: replaces,
			Maintainers: []csvv1alpha1.Maintainer{
				csvv1alpha1.Maintainer{
					Name:  "OpenStack K8s",
					Email: "openstack-discuss@lists.openstack.org",
				},
			},
			Maturity: "alpha",
			Provider: csvv1alpha1.AppLink{
				Name: "OpenStack K8s",
			},
			Links: []csvv1alpha1.AppLink{
				csvv1alpha1.AppLink{
					Name: "Source Code",
					URL:  "https://github.com/openstack-k8s-operators/openstack-cluster-operator",
				},
			},
			Icon: []csvv1alpha1.Icon{
				csvv1alpha1.Icon{
					MediaType: "image/png",
					Data:      "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AABArAEAACDWAAAAxBoAAABiDQAAALEGAAAg1gAAABBrAAAAYg0AAACxBgAAINYAAAAQawAAAIg1AAAAsQYAAIBYAwAAEGsAAACINQAAALEGAACAWAMAAECsAQAAiDUAAADEGgAAQGtc+s+/+WQlyZqhoGF2O/3qwDAAMG+Dbm/DKNAwB51+tWsY6u9yCbXKUNAwvSTbhgGAGvA+iqZ5l+QTw1B/lkECAACINQAAAMQaAACAWAMAAECsAQAAiDUAAADEGgAAgFgDAABArAEAACDWAAAAxBoAAABiDQAAQKwBAAAg1gAAAMQaAAAAYg0AAACxBgAAINYAAAAQawAAAGINAAAAsQYAACDWAAAAEGsAAACINQAAALEGAACAWAMAABBrAAAAiDUAAACxBgAAgFgDAABArAEAAIg1AAAAxBoAAIBYAwAAQKwBAACINQAAAMQaAAAAYg0AAECsAQAAINYAAADEGgAAAGINAABArAEAACDWAAAAEGsAAABiDQAAALEGAAAg1gAAABBrAAAAYg0AAACxBgAAgFgDAAAQawAAAIg1AAAAsQYAAIBYAwAAEGsAAACINQAAAMQaAACAWAMAAECsAQAAiDUAAADEGgAAgFgDAABArAEAACDWAAAAxBoAAABiDQAAQKwBAAAg1gAAAMQaAAAAYg0AAACxBgAAINYAAAAQawAAAGINAAAAsQYAACDWAAAAEGsAAACINQAAALEGAACAWAMAABBrAAAAiDUAAACxBgAAgFgDAABArAEAAIg1AAAAxBoAAIBYAwAAQKwBAACINQAAAMQaAACAWAMAAECsAQAAINYAAADEGgAAAGINAABArAEAACDWAAAAxBoAAABiDQAAALEGAAAg1gAAABBrAAAAYg0AAACxBgAAINYAAAAQawAAAIg1AAAAsQYAAIBYAwAAEGsAAACINQAAALEGAACAWAMAAECsAQAAiDUAAADEGgAAgFgDAABArAEAAIg1AAAAxBoAAABiDQAAQKwBAAAg1gAAAMQaAAAAYg0AAECsAQAAINYAAAAQawAAAGINAAAAsQYAACDWAAAAEGsAAABiDQAAALEGAACAWAMAABBrAAAAiDUAAACxBgAAgFgDAAAQawAAAIg1AAAAxBoAAIBYAwAAQKwBAACINQAAAMQaAACAWAMAAECsAQAAINYAAADEGgAAAGINAABA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",
				},
			},
			Labels: map[string]string{
				"alm-owner-openstack-cluster-operator": "openstack-cluster-operator",
				"operated-by":                          "openstack-cluster",
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"alm-owner-openstack-cluster-operator": "openstack-cluster-operator",
					"operated-by":                          "openstack-cluster-operator",
				},
			},
			InstallModes: []csvv1alpha1.InstallMode{
				csvv1alpha1.InstallMode{
					Type:      csvv1alpha1.InstallModeTypeOwnNamespace,
					Supported: true,
				},
				csvv1alpha1.InstallMode{
					Type:      csvv1alpha1.InstallModeTypeSingleNamespace,
					Supported: true,
				},
				csvv1alpha1.InstallMode{
					Type:      csvv1alpha1.InstallModeTypeMultiNamespace,
					Supported: false,
				},
				csvv1alpha1.InstallMode{
					Type:      csvv1alpha1.InstallModeTypeAllNamespaces,
					Supported: false,
				},
			},
			InstallStrategy: csvv1alpha1.NamedInstallStrategy{},
			CustomResourceDefinitions: csvv1alpha1.CustomResourceDefinitions{
				Owned: []csvv1alpha1.CRDDescription{
					csvv1alpha1.CRDDescription{
						Name:        "openstackclusters.openstackcluster.openstack.org",
						Version:     "v1",
						Kind:        "OpenStackCluster",
						DisplayName: crdDisplay + " Deployment",
						Description: "Represents the deployment of " + crdDisplay,
					},
				},
				Required: []csvv1alpha1.CRDDescription{},
			},
		},
	}
}

func int32Ptr(i int32) *int32 {
	return &i
}
