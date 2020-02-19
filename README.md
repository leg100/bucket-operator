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
