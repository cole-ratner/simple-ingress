package pod

import (
	corev1 "k8s.io/api/core/v1"
	networkingv1 "simpleingress/api/v1"
)

// New returns a new instance of a corev1.Pod that is based on a kolamiti92/simpleproxy container
func New(cr *networkingv1.SimpleIngress) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: cr.ObjectMeta,
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "SimpleIngressReverseProxy-{{.cr.Host}}",
					Image: "kolamiti92/simpleproxy",
					Ports: []corev1.ContainerPort{
						{
							Name:          "http",
							ContainerPort: 8080,
							Protocol:      "TCP",
						},
					},
				},
			},
			RestartPolicy: corev1.RestartPolicyOnFailure,
		},
	}
}
