package bucket

import (
	"context"

	"cloud.google.com/go/storage"
	goalspikev1alpha1 "github.com/leg100/bucket-operator/pkg/apis/goalspike/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_bucket")

// Add creates a new Bucket Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	reconciler, err := newReconciler(mgr)
	if err != nil {
		return err
	}
	return add(mgr, reconciler)
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) (reconcile.Reconciler, error) {
	storageClient, err := storage.NewClient(context.TODO())
	if err != nil {
		return nil, err
	}
	return &ReconcileBucket{storageClient: storageClient, client: mgr.GetClient(), scheme: mgr.GetScheme()}, nil
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("bucket-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Bucket
	err = c.Watch(&source.Kind{Type: &goalspikev1alpha1.Bucket{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Bucket
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &goalspikev1alpha1.Bucket{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileBucket implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileBucket{}

// ReconcileBucket reconciles a Bucket object
type ReconcileBucket struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client        client.Client
	storageClient *storage.Client
	scheme        *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Bucket object and makes changes based on the state read
// and what is in the Bucket.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileBucket) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Bucket")

	// Fetch the Bucket instance
	instance := &goalspikev1alpha1.Bucket{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	name := instance.Spec.Name
	project := instance.Spec.Project

	bh := r.storageClient.Bucket(name)
	_, err = bh.Attrs(context.TODO())
	if err != nil {
		if err == storage.ErrBucketNotExist {
			reqLogger.Info("Bucket not found, creating bucket", "Bucket.Spec.Name", name)
			if err = bh.Create(context.TODO(), project, nil); err != nil {
				return reconcile.Result{}, err
			} else {
				return reconcile.Result{Requeue: true}, nil
			}
		} else {
			return reconcile.Result{}, err
		}
	}

	// Bucket already exists - don't requeue
	reqLogger.Info("Skip reconcile: Bucket already exists", "Bucket.Spec.Name", name)
	return reconcile.Result{}, nil
}
