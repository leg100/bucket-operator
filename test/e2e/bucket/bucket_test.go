package e2e

import (
	"testing"
	"time"

	goctx "context"

	"github.com/leg100/bucket-operator/pkg/apis"
	cachev1alpha1 "github.com/leg100/bucket-operator/pkg/apis/goalspike/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
)

var (
	retryInterval        = time.Second * 5
	timeout              = time.Second * 60
	cleanupRetryInterval = time.Second * 1
	cleanupTimeout       = time.Second * 5
)

func TestBucket(t *testing.T) {
	bucketList := &cachev1alpha1.BucketList{}
	err := framework.AddToFrameworkScheme(apis.AddToScheme, bucketList)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}

	t.Run("Bucket1", CreateBucket)
}

func CreateBucket(t *testing.T) {
	ctx := framework.NewTestCtx(t)
	defer ctx.Cleanup()

	err := ctx.InitializeClusterResources(&framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		t.Fatalf("failed to initialize cluster resources: %v", err)
	}

	// get namespace
	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Fatal(err)
	}
	// get global framework variables
	f := framework.Global
	// wait for bucket-operator to be ready
	err = e2eutil.WaitForOperatorDeployment(t, f.KubeClient, namespace, "bucket-operator", 1, time.Second*5, time.Second*30)
	if err != nil {
		t.Fatal(err)
	}

	// create bucket custom resource
	exampleBucket := &cachev1alpha1.Bucket{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-bucket",
			Namespace: namespace,
		},
		Spec: cachev1alpha1.BucketSpec{
			Name:    "automatize-my-new-bucket",
			Project: "automatize-admin",
		},
	}
	err = f.Client.Create(goctx.TODO(), exampleBucket, &framework.CleanupOptions{TestContext: ctx, Timeout: time.Second * 5, RetryInterval: time.Second * 1})
	if err != nil {
		t.Fatal(err)
	}
}
