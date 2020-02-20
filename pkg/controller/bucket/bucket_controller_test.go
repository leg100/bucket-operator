package bucket

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/fsouza/fake-gcs-server/fakestorage"
	//TODO: rename
	cachev1alpha1 "github.com/leg100/bucket-operator/pkg/apis/goalspike/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

func setupFakeStorage(t *testing.T) *fakestorage.Server {
	opts := fakestorage.Options{
		StorageRoot: "testdata",
	}
	if testing.Verbose() {
		opts.Writer = os.Stdout
	}
	server, err := fakestorage.NewServerWithOptions(opts)

	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("server started at %s", server.URL())

	return server
}

// TestBucketController runs ReconcileBucket.Reconcile() against a
// fake client that tracks a Bucket object.
func TestBucketController(t *testing.T) {
	// Set the logger to development mode for verbose logs.
	logf.SetLogger(logf.ZapLogger(true))

	// setup fake GCS storage backend
	server := setupFakeStorage(t)
	storageClient := server.Client()

	var (
		name              = "bucket-operator"
		namespace         = "bucket"
		bucketName string = "my-new-bucket"
	)

	// A Bucket resource with metadata and spec.
	bucket := &cachev1alpha1.Bucket{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: cachev1alpha1.BucketSpec{
			Name: bucketName,
		},
	}
	// Objects to track in the fake client.
	objs := []runtime.Object{
		bucket,
	}

	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(cachev1alpha1.SchemeGroupVersion, bucket)
	// Create a fake client to mock API calls.
	cl := fake.NewFakeClient(objs...)
	// Create a ReconcileBucket object with the scheme and fake client.
	r := &ReconcileBucket{storageClient: storageClient, client: cl, scheme: s}

	// Mock request to simulate Reconcile() being called on an event for a
	// watched resource .
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		},
	}
	res, err := r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
	// Check the result of reconciliation to make sure it has the desired state.
	if !res.Requeue {
		t.Error("reconcile did not requeue request as expected")
	}

	bh := storageClient.Bucket(bucketName)
	// Next check if the bucket exists
	if _, err := bh.Attrs(context.TODO()); err != nil {
		t.Errorf("expected bucket %s to exist", bucketName)
	}

	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
	// Check the result of reconciliation to make sure it has the desired state.
	if res.Requeue {
		t.Error("reconcile requeue which is not expected")
	}

	if err = os.Remove("./testdata/my-new-bucket"); err != nil {
		t.Fatalf("couldn't delete bucket %s", bucketName)
	}

	// Get the updated Bucket object.
	bucket = &cachev1alpha1.Bucket{}
	err = r.client.Get(context.TODO(), req.NamespacedName, bucket)
	if err != nil {
		t.Errorf("get bucket: (%v)", err)
	}

	// Ensure Reconcile() updated the Bucket's Status as expected.
	// nodes := bucket.Status.Nodes
	// if !reflect.DeepEqual(podNames, nodes) {
	// 	t.Errorf("pod names %v did not match expected %v", nodes, podNames)
	// }

	server.Stop()
}
