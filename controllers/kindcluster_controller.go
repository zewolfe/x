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
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/pkg/errors"
	infrastructurev1alpha4 "github.com/zewolfe/cluster-api-provider-kind/api/v1alpha4"
)

// KINDClusterReconciler reconciles a KINDCluster object
type KINDClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=kindclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=kindclusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=kindclusters/finalizers,verbs=update
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the KINDCluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *KINDClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	kindCluster := &infrastructurev1alpha4.KINDCluster{}

	if err := r.Get(ctx, req.NamespacedName, kindCluster); err != nil {

		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	log = log.WithValues("cluster", kindCluster.Name)
	helper, err := patch.NewHelper(kindCluster, r.Client)
	if err != nil {
		return reconcile.Result{}, errors.Wrap(err, "failed to init patch helper")
	}

	if err := helper.Patch(ctx, kindCluster); err != nil {
		return ctrl.Result{}, errors.Wrapf(err, "couldn't patch cluster %q", kindCluster.Name)
	}

	kindCluster.Status.Ready = true

	// Fetch the Cluster.
	cluster, err := util.GetOwnerCluster(ctx, r.Client, kindCluster.ObjectMeta)

	fmt.Println(kindCluster.ObjectMeta, "After cluster has been fetched")
	fmt.Println("/n")
	// fmt.Println(cluster, "After cluster has been fetched")
	// log.Info("After cluster has been fetched")
	log.Info(kindCluster.ObjectMeta.ClusterName)
	if err != nil {
		return reconcile.Result{}, err
	}

	if cluster == nil {
		log.Info("Cluster Controller has not yet set OwnerRef")
		return reconcile.Result{}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KINDClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrastructurev1alpha4.KINDCluster{}).
		// Watches(
		// 	&source.Kind{Type: &clusterv1.Cluster{}},
		// 	&handler.EnqueueRequestsFromMapFunc{
		// 		ToRequests: util.ClusterToInfrastructureMapFunc(infrastructurev1alpha4.GroupVersion.WithKind("KINDCluster")),
		// 	}).
		Complete(r)
}
