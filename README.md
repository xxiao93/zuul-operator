# zuul-operator

## Prerequisites
* go 1.12.9 linux/amd64
* operator-sdk v0.11.0

## Quick Start

### install operator
https://github.com/operator-framework/operator-sdk

### init operator project
```apple js
$ mkdir -p $GOPATH/src/github.com/example-inc/
$ cd $GOPATH/src/github.com/example-inc/
$ export GO111MODULE=on
$ operator-sdk new zuul-operator
$ cd zuul-operator
```

then git clone zuul-operator to $GOPATH/src/github.com/example-inc/ and run next cmd
```apple js
kubectl create -f deploy/crds/cache_v1alpha1_zuul_crd.yaml
kubectl create -f deploy/crds/cache_v1alpha1_zuul_cr.yaml
operator-sdk up local --namespace=default
```