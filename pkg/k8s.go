package k8s

import (
	"encoding/json"
	operatorv1 "juggernaut/api/v1"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"os"
)

const (
	generatedFromAnnotation = "oceanoperator.com/generated-from"
)

var (
	//juggernaut默认使用的镜像
	defaultImage = "juggernaut:v1.3"
	// 默认副本数
	defaultReplicas      = int32(1)
	defaultContainerPort = int32(8080)
)

// NewDeployment creates a deployment for a given Nginx resource.
func NewDeployment(juggernaut *operatorv1.Juggernaut) (*appv1.Deployment, error) {
	deployment := appv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      juggernaut.Name + "-deployment",
			Namespace: juggernaut.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(juggernaut, schema.GroupVersionKind{
					Group:   operatorv1.GroupVersion.Group,
					Version: operatorv1.GroupVersion.Version,
					Kind:    "Juggernaut",
				}),
			},
			Labels: LabelsForJuggernaut(juggernaut.Name),
		},
		Spec: appv1.DeploymentSpec{
			Replicas: &defaultReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: LabelsForJuggernaut(juggernaut.Name),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: LabelsForJuggernaut(juggernaut.Name),
				},
				Spec: corev1.PodSpec{
					EnableServiceLinks: func(b bool) *bool { return &b }(false),
					Containers: []corev1.Container{
						{
							Name:      "juggernaut",
							Image:     juggernaut.Spec.Image,
							Resources: juggernaut.Spec.Resources,
							Ports: []corev1.ContainerPort{{
								ContainerPort: defaultContainerPort,
								Protocol:      "TCP",
							}},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      juggernaut.Name + "-configmap",
									MountPath: "/JUGGERNAUT/config",
									ReadOnly:  false,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: juggernaut.Name + "-configmap",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: juggernaut.Name + "-configmap",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// This is done on the last step because n.Spec may have mutated during these methods
	if err := SetJuggernautSpec(&deployment.ObjectMeta, juggernaut.Spec); err != nil {
		return nil, err
	}

	return &deployment, nil
}

// LabelsForJuggernaut returns the labels for a juggernaut CR with the given name
func LabelsForJuggernaut(name string) map[string]string {
	return map[string]string{
		"oceanoperator.com/resource-name": name,
		"oceanoperator.com/app":           "juggernaut",
	}
}

// SetJuggernautSpec sets the juggernaut spec into the object annotation to be later extracted
func SetJuggernautSpec(o *metav1.ObjectMeta, spec operatorv1.JuggernautSpec) error {
	if o.Annotations == nil {
		o.Annotations = make(map[string]string)
	}
	origSpec, err := json.Marshal(spec)
	if err != nil {
		return err
	}
	o.Annotations[generatedFromAnnotation] = string(origSpec)
	return nil
}

// NewService 实现一个根据juggernaut对象生成的service对象
func NewService(juggernaut *operatorv1.Juggernaut) *corev1.Service {
	service := corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      juggernaut.Name + "-service",
			Namespace: juggernaut.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(juggernaut, schema.GroupVersionKind{
					Group:   operatorv1.GroupVersion.Group,
					Version: operatorv1.GroupVersion.Version,
					Kind:    "Juggernaut",
				}),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: LabelsForJuggernaut(juggernaut.Name),
			Type:     juggernaut.Spec.Service.Type,
			Ports: []corev1.ServicePort{
				{
					Port:       int32(8080),
					Protocol:   "TCP",
					TargetPort: intstr.FromInt32(defaultContainerPort),
				},
			},
		},
	}
	return &service
}

//  实现一个根据juggernaut对象生成的configmap对象

func NewConfigmap(juggernaut *operatorv1.Juggernaut) *corev1.ConfigMap {
	//默认对象
	configmap := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      juggernaut.Name + "-configmap",
			Namespace: juggernaut.Namespace,
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(juggernaut, schema.GroupVersionKind{
				Group:   operatorv1.GroupVersion.Group,
				Version: operatorv1.GroupVersion.Version,
				Kind:    "Juggernaut",
			})},
		},
		Data: defaultConfigMapData,
	}
	return &configmap
}

func readDefaultConfig() (a string) {
	file, err := os.ReadFile("./config.yaml")
	if err != nil {
		return
	}
	a = string(file)
	return a
}

var (
	defaultConfigMapData = map[string]string{"config.yaml": readDefaultConfig()}
)
