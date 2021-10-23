/*
Copyright 2021.

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

package controllers

import (
	"context"
	"encoding/json"
	"reflect"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	datameshv1 "github.com/odemeulder/datamesh/api/v1"
)

// DataProductReconciler reconciles a DataProduct object
type DataProductReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=datamesh.demeulder.us,resources=dataproducts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=datamesh.demeulder.us,resources=dataproducts/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=datamesh.demeulder.us,resources=dataproducts/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;
//
// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DataProduct object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.9.2/pkg/reconcile
func (r *DataProductReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.Log.WithName("reconciler")

	// Fetch the Memchached instance
	dp := &datameshv1.DataProduct{}
	err := r.Get(ctx, req.NamespacedName, dp)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
	}

	// Check if configmaps already exists, if not create a new configmap
	cms := r.configMapsForDeployment(ctx, dp)
	for _, cm := range cms {
		foundCm := &corev1.ConfigMap{}
		err = r.Get(ctx, types.NamespacedName{Name: cm.Name, Namespace: cm.Namespace}, foundCm)
		if err != nil {
			if errors.IsNotFound(err) {
				if err = r.Create(ctx, &cm); err != nil {
					log.Info("Could not create configmap")
					return ctrl.Result{}, err
				}
			}
		} else {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// foundCm := &corev1.ConfigMap{}
	// err = r.Get(ctx, types.NamespacedName{Name: dp.Name, Namespace: dp.Namespace}, foundCm)
	// if err != nil {
	// 	if errors.IsNotFound(err) {
	// 		cms := r.configMapsForDeployment(ctx, dp)
	// 		for _, cm := range cms {
	// 			if err = r.Create(ctx, &cm); err != nil {
	// 				log.Info("Could not create configmap")
	// 				return ctrl.Result{}, err
	// 			}
	// 		}
	// 		return ctrl.Result{Requeue: true}, nil
	// 	} else {
	// 		return ctrl.Result{}, err
	// 	}
	// }

	// Check if the deployment already exists, if not create a new deployment
	found := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: dp.Name, Namespace: dp.Namespace}, found)
	if err != nil {
		if errors.IsNotFound(err) {
			// Define and create a new deployment
			dep := r.deploymentForMemCached(dp)
			if err = r.Create(ctx, dep); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil
		} else {
			return ctrl.Result{}, err
		}
	}
	// Ensure the deployment size is the same as the spec
	size := dp.Spec.Size
	if *found.Spec.Replicas != size {
		found.Spec.Replicas = &size
		if r.Update(ctx, found); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// Update the Memcached status with the pod names
	// List the pods for this CR's deployment
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(dp.Namespace),
		client.MatchingLabels(labelsForApp(dp.Name)),
	}
	if err = r.List(ctx, podList, listOpts...); err != nil {
		return ctrl.Result{}, err
	}
	podNames := getPodNames(podList.Items)
	if !reflect.DeepEqual(podNames, dp.Status.Nodes) {
		dp.Status.Nodes = podNames
		if err := r.Status().Update(ctx, dp); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil

}

func getPodNames(pods []corev1.Pod) []string {
	var s []string
	for _, v := range pods {
		s = append(s, v.Name)
	}
	return s
}

func (r *DataProductReconciler) configMapsForDeployment(ctx context.Context, dp *datameshv1.DataProduct) []corev1.ConfigMap {
	log := ctrl.Log.WithName("reconciler")

	inputPorts := dp.Spec.InputPorts
	var err error
	var cms []corev1.ConfigMap
	for _, v := range inputPorts {
		ip := &datameshv1.InputPortS3{}
		err = r.Get(ctx, types.NamespacedName{Name: v, Namespace: dp.Namespace}, ip)
		if err != nil {
			log.Info("Input port not found", "port", v)
			continue
		}
		ipJson, err := json.Marshal(ip)
		if err != nil {
			log.Info("Error unmarshalling Json for input port", "port", v)
			continue
		}
		log.Info("Unmarshalled:", "blob", string(ipJson))
		cm := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      v,
				Namespace: dp.Namespace,
			},
			Data: map[string]string{
				v: string(ipJson),
			},
		}
		cms = append(cms, *cm)
	}

	return cms
}

// deploymentForMemcached returns a Deployment object for data from m
func (r *DataProductReconciler) deploymentForMemCached(dp *datameshv1.DataProduct) *appsv1.Deployment {
	lbls := labelsForApp(dp.Name)
	replicas := dp.Spec.Size

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dp.Name,
			Namespace: dp.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: lbls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: lbls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:   dp.Spec.Image,
						Name:    "dataproduct",
						Command: []string{"sleep", "10m"},
						Ports: []corev1.ContainerPort{{
							ContainerPort: 10000,
							Name:          "dataproduct",
						}},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "config-volume",
							MountPath: "/etc/config/odm.conf",
							SubPath:   "odm.conf",
						}},
					}},
					Volumes: []corev1.Volume{{
						Name: "config-volume",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: dp.Name,
								},
							},
						},
					}},
				},
			},
		},
	}

	controllerutil.SetControllerReference(dp, dep, r.Scheme)
	return dep
}

// labelsForApp creates a simple set of labels for Dataproduct.
func labelsForApp(name string) map[string]string {
	return map[string]string{"cr_name": name}
}

// SetupWithManager sets up the controller with the Manager.
func (r *DataProductReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&datameshv1.DataProduct{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}
