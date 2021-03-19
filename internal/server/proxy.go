package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	networkingv1 "simpleingress/api/v1"
)

func main() {

	s := &http.Server{}

	log.Fatal(s.ListenAndServe())

}


func ConfigProxy(cr *networkingv1.SimpleIngress) {

}


func NewPod(cr *networkingv1.SimpleIngress) {
	return &corev1.Pod{
		ObjectMeta: cr.ObjectMeta,
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name: "SimpleIngressReverseProxy",
					Image: "gravityloop/simpleingress-reverseproxy",
				},
			},
			RestartPolicy: corev1.RestartPolicyOnFailure,
		},
	}
}