# bucket-operator

[![leg100](https://circleci.com/gh/leg100/bucket-operator.svg?style=svg)](https://circleci.com/gh/leg100/bucket-operator)

A kubernetes operator for GCS buckets.

## install

Deploy the `Bucket` CRD:

```bash
kubectl apply -f https://raw.githubusercontent.com/leg100/bucket-operator/master/deploy/crds/goalspike.com_buckets_crd.yaml
```

Create and deploy a secret, giving it the path to a GCP service account key (the service account needs to have the permission to create buckets):

```bash
kubectl create secret generic bucket-operator \
  --from-file=key.json=<path_to_key_file>
```

Deploy the controller including its service account and RBAC resources:

```bash
kubectl apply -f https://raw.githubusercontent.com/leg100/bucket-operator/master/deploy/service_account.yaml
kubectl apply -f https://raw.githubusercontent.com/leg100/bucket-operator/master/deploy/role.yaml
kubectl apply -f https://raw.githubusercontent.com/leg100/bucket-operator/master/deploy/role_binding.yaml
kubectl apply -f https://raw.githubusercontent.com/leg100/bucket-operator/master/deploy/operator.yaml
```

## use

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
kubectl logs -f bucket-operator-7fb84b479c-cxql5
...
{"level":"info","ts":1582152903.7753346,"logger":"controller_bucket","msg":"Bucket not found, creating bucket","Request.Namespace":"default","Request.Name":"example-bucket","Bucket.Spec.Name":"automatize-test-bucket"}
{"level":"info","ts":1582152905.1113892,"logger":"controller_bucket","msg":"Reconciling Bucket","Request.Namespace":"default","Request.Name":"example-bucket"}
{"level":"info","ts":1582152905.4145389,"logger":"controller_bucket","msg":"Skip reconcile: Bucket already exists","Request.Namespace":"default","Request.Name":"example-bucket","Bucket.Spec.Name":"automatize-test-bucket"}
```

The logs confirm:

* if the bucket doesn't exist, it creates one
* it'll then notice it now exists
