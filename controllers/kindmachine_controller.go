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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrastructurev1alpha4 "github.com/zewolfe/cluster-api-provider-kind/api/v1alpha4"
)

// KINDMachineReconciler reconciles a KINDMachine object
type KINDMachineReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=kindmachines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=kindmachines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=kindmachines/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the KINDMachine object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *KINDMachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// log := log.FromContext(ctx)

	// kindMachine := &infrastructurev1alpha4.KINDMachine{}
	// if err := r.Client.Get(ctx, req.NamespacedName, kindMachine); err != nil {
	// 	if apierrors.IsNotFound(err) {
	// 		return ctrl.Result{}, nil
	// 	}
	// 	return ctrl.Result{}, err
	// }

	// // Fetch the Machine.
	// machine, err := util.GetOwnerMachine(ctx, r.Client, kindMachine.ObjectMeta)
	// if err != nil {
	// 	return ctrl.Result{}, err
	// }
	// if machine == nil {
	// 	log.Info("Waiting for Machine Controller to set OwnerRef on DockerMachine")
	// 	return ctrl.Result{}, nil
	// }

	// log = log.WithValues("machine", machine.Name)

	// cluster, err := util.GetClusterFromMetadata(ctx, r.Client, machine.ObjectMeta)
	// if err != nil {
	// 	log.Info("KindMachine owner Machine is missing cluster label or cluster does not exist")
	// 	return ctrl.Result{}, err
	// }
	// if cluster == nil {
	// 	log.Info(fmt.Sprintf("Please associate this machine with a cluster using the label %s: <name of cluster>", clusterv1.ClusterLabelName))
	// 	return ctrl.Result{}, nil
	// }

	// log = log.WithValues("cluster", cluster.Name)

	// // Fetch the Docker Cluster.
	// dockerCluster := &infrastructurev1alpha4.KINDCluster{}
	// dockerClusterName := client.ObjectKey{
	// 	Namespace: kindMachine.Namespace,
	// 	Name:      cluster.Spec.InfrastructureRef.Name,
	// }
	// if err := r.Client.Get(ctx, dockerClusterName, dockerCluster); err != nil {
	// 	log.Info("KindCluster is not available yet")
	// 	return ctrl.Result{}, nil
	// }

	// log = log.WithValues("kind-cluster", dockerCluster.Name)

	// // Initialize the patch helper
	// patchHelper, err := patch.NewHelper(kindMachine, r.Client)
	// if err != nil {
	// 	return ctrl.Result{}, err
	// }
	// // Always attempt to Patch the DockerMachine object and status after each reconciliation.
	// defer func() {
	// 	if err := patchHelper.Patch(ctx, kindMachine); err != nil {
	// 		log.Error(err, "failed to patch DockerMachine")
	// 	}
	// }()

	// kindMachine.Status.Ready = true

	// conditions.MarkTrue(kindMachine, clusterv1.BootstrapReadyCondition)
	// your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KINDMachineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrastructurev1alpha4.KINDMachine{}).
		Complete(r)
}
