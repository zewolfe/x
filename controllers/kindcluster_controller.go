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
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	infrastructurev1alpha4 "github.com/zewolfe/cluster-api-provider-kind/api/v1alpha4"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha4"
)

// KINDClusterReconciler reconciles a KINDCluster object
type KINDClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

type clusterReconcileContext struct {
	ctx            context.Context
	kindCluster    *infrastructurev1alpha4.KINDCluster
	patchHelper    *patch.Helper
	cluster        *clusterv1.Cluster
	log            logr.Logger
	client         client.Client
	namespacedName types.NamespacedName
}

func (r *KINDClusterReconciler) newReconcileContext(ctx context.Context, req ctrl.Request) (*clusterReconcileContext, error) {
	log := log.FromContext(ctx)
	crc := &clusterReconcileContext{
		log:            log.WithValues("KINDCluster", req.NamespacedName),
		ctx:            ctx,
		kindCluster:    &infrastructurev1alpha4.KINDCluster{},
		client:         r.Client,
		namespacedName: req.NamespacedName,
	}

	if err := crc.client.Get(crc.ctx, crc.namespacedName, crc.kindCluster); err != nil {
		if apierrors.IsNotFound(err) {
			crc.log.Info("KINDCluster Object not found")
			return nil, nil
		}

		return nil, err
	}

	helper, err := patch.NewHelper(crc.kindCluster, crc.client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}

	crc.patchHelper = helper

	cluster, err := util.GetOwnerCluster(crc.ctx, crc.client, crc.kindCluster.ObjectMeta)
	// cluster, err := util.GetClusterFromMetadata(crc.ctx, crc.client, crc.kindCluster.ObjectMeta)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			crc.log.Info("Error getting owner cluster")
			return nil, err
		}
	}

	// x, err := util.GetClusterFromMetadata()

	// if cluster == nil {
	// 	crc.log.Info("Owner Cluster not set. Requeue")
	// 	return reconcile.Result{}, nil
	// }

	crc.cluster = cluster

	return crc, nil
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
func (r *KINDClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, rerr error) {
	log := log.FromContext(ctx)

	log.Info("Starting new Reconcile context")

	crc, err := r.newReconcileContext(ctx, req)
	if err != nil {
		fmt.Println("Error creating reconciliation context: %w", err)
		return ctrl.Result{}, err
	}

	if crc == nil {
		return ctrl.Result{}, nil
	}

	// if crc.cluster == nil {
	// 	return ctrl.Result{}, nil
	// }

	crc.kindCluster.Status.Ready = true
	conditions.MarkTrue(crc.kindCluster, clusterv1.ReadyCondition)
	// conditions.MarkTrue(crc.kindCluster, clusterv1.ConditionType(clusterv1.ClusterPhaseProvisioned))

	// Always attempt to Patch the DockerCluster object and status after each reconciliation.
	defer func() {
		crc.log.Info("Setting cluster status to ready")
		if err := crc.patchHelper.Patch(crc.ctx, crc.kindCluster, patch.WithOwnedConditions{
			Conditions: []clusterv1.ConditionType{
				clusterv1.ReadyCondition,
				clusterv1.ConditionType(clusterv1.ClusterPhaseProvisioned),
			},
		}); err != nil {
			fmt.Println("patching cluster object: %w", err)

			// return ctrl.Result{}, nil
		}
	}()
	// crc.log.Info("Setting cluster status to ready")
	// if err := crc.patchHelper.Patch(crc.ctx, crc.kindCluster); err != nil {
	// 	fmt.Println("patching cluster object: %w", err)

	// 	return ctrl.Result{}, nil
	// }

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KINDClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	c, err := ctrl.NewControllerManagedBy(mgr).
		For(&infrastructurev1alpha4.KINDCluster{}).
		Build(r)

	if err != nil {
		return err
	}

	return c.Watch(
		&source.Kind{Type: &clusterv1.Cluster{}},
		handler.EnqueueRequestsFromMapFunc(util.ClusterToInfrastructureMapFunc(infrastructurev1alpha4.GroupVersion.WithKind("KINDCluster"))),
	)
}
