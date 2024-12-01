package k8s

import (
	"encoding/json"
	operatorv1 "juggernaut/api/v1"
	"strings"

	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
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
			Name:      juggernaut.Name,
			Namespace: juggernaut.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(juggernaut, schema.GroupVersionKind{
					Group:   operatorv1.GroupVersion.Group,
					Version: operatorv1.GroupVersion.Version,
					Kind:    "Juggernaut",
				}),
			},
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
					Containers: append([]corev1.Container{
						{
							Name:      "juggernaut",
							Image:     defaultImage,
							Resources: juggernaut.Spec.Resources,
							Ports: []corev1.ContainerPort{{
								ContainerPort: defaultContainerPort,
								Protocol:      "TCP",
							}},
						},
					}),
					Volumes: []corev1.Volume{
						{
							Name: transName(juggernaut.Spec.Config.Overwrite.Name),
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: transName(juggernaut.Spec.Config.Overwrite.Name),
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

func transName(name operatorv1.NamespacedName) string {
	a := strings.Split(string(name), "/")
	return a[1]
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

}

//  实现一个根据juggernaut对象生成的service对象

func NewConfigmap(juggernaut *operatorv1.Juggernaut) *corev1.ConfigMap {

}
