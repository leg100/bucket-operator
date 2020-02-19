# bucket-operator

[![leg100](https://circleci.com/gh/leg100/bucket-operator.svg?style=svg)](https://circleci.com/gh/leg100/bucket-operator)

A kubernetes operator for GCS buckets.

## install

Deploy the `Bucket` CRD:

```bash
kubectl apply -f https://raw.githubusercontent.com/leg100/bucket-operator/master/deploy/crds/goalspike.com_buckets_crd.yaml
```

Deploy the controller:

```bash
kubectl apply -f https://raw.githubusercontent.com/leg100/bucket-operator/master/deploy/operator.yaml
```

Create and deploy a `Bucket` custom resource:

```bash
cat > bucket.yaml <EOF
apiVersion: goalspike.com/v1alpha1
kind: Bucket
metadata:
  name: example-bucket
spec:
  name: <my-gcs-bucket-name>
  project: <my-gcp-project>
EOF
 
kubectl apply -f bucket.yaml
```

Check the controller's logs:

```bash
~/co/bucket-operator(master) $ kubectl logs -f bucket-operator-7fb84b479c-cxql5
...
{"level":"info","ts":1582152903.7753346,"logger":"controller_bucket","msg":"Bucket not found, creating bucket","Request.Namespace":"default","Request.Name":"example-bucket","Bucket.Spec.Name":"automatize-test-bucket"}
{"level":"info","ts":1582152905.1113892,"logger":"controller_bucket","msg":"Reconciling Bucket","Request.Namespace":"default","Request.Name":"example-bucket"}
{"level":"info","ts":1582152905.4145389,"logger":"controller_bucket","msg":"Skip reconcile: Bucket already exists","Request.Namespace":"default","Request.Name":"example-bucket","Bucket.Spec.Name":"automatize-test-bucket"}
```

