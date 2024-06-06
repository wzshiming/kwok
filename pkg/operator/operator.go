package operator

import (
	"context"

	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/kwok/pkg/client/operator/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/kwok/pkg/apis/operator/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type Operator struct {
	typedClient kubernetes.Interface
	client      versioned.Interface
}

func NewOperator() *Operator {
	return nil
}

func (o *Operator) Start(ctx context.Context) error {
	o.client.OperatorV1alpha1().Kwoks(metav1.NamespaceAll)
	return nil
}

func crToDeployment(k *v1alpha1.Kwok) (*appsv1.Deployment, error) {
	d := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        k.Name,
			Namespace:   k.Namespace,
			Labels:      k.Labels,
			Annotations: k.Annotations,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: k.APIVersion,
					Kind:       k.Kind,
					Name:       k.Name,
					UID:        k.UID,
				},
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: k.Spec.Selector,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: k.Spec.Selector.MatchLabels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: k.Spec.Image,
							Env: []corev1.EnvVar{
								{
									Name: "POD_IP",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "status.podIP",
										},
									},
								},
								{
									Name: "HOST_IP",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "status.hostIP",
										},
									},
								},
							},
							Args: []string{
								"--manage-all-nodes=false",
								"--manage-nodes-with-annotation-selector=kwok.x-k8s.io/node=fake",
								"--manage-nodes-with-label-selector=",
								"--manage-single-node=",
								"--node-ip=$(POD_IP)",
								"--node-port=10247",
								"--cidr=10.0.0.1/24",
								"--node-lease-duration-seconds=40",
								"--enable-crds=Stage",
								"--enable-crds=Metric",
								"--enable-crds=Attach",
								"--enable-crds=ClusterAttach",
								"--enable-crds=Exec",
								"--enable-crds=ClusterExec",
								"--enable-crds=Logs",
								"--enable-crds=ClusterLogs",
								"--enable-crds=PortForward",
								"--enable-crds=ClusterPortForward",
								"--enable-crds=ResourceUsage",
								"--enable-crds=ClusterResourceUsage",
							},
							LivenessProbe: &corev1.Probe{
								FailureThreshold:    10,
								InitialDelaySeconds: 30,
								PeriodSeconds:       60,
								TimeoutSeconds:      10,
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path:   "/healthz",
										Port:   intstr.FromInt32(10247),
										Scheme: corev1.URISchemeHTTPS,
									},
								},
							},
							ReadinessProbe: &corev1.Probe{
								FailureThreshold:    5,
								InitialDelaySeconds: 2,
								PeriodSeconds:       20,
								TimeoutSeconds:      2,
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path:   "/healthz",
										Port:   intstr.FromInt32(10247),
										Scheme: corev1.URISchemeHTTPS,
									},
								},
							},
							StartupProbe: &corev1.Probe{
								FailureThreshold:    3,
								InitialDelaySeconds: 2,
								PeriodSeconds:       10,
								TimeoutSeconds:      2,
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path:   "/healthz",
										Port:   intstr.FromInt32(10247),
										Scheme: corev1.URISchemeHTTPS,
									},
								},
							},
						},
					},

					RestartPolicy:      corev1.RestartPolicyAlways,
					ServiceAccountName: "kwok-controller",
				},
			},
		},
	}

	return &d, nil
}
