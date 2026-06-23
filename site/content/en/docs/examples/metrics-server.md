---
title: "Metrics Server"
---

# Metrics Server

More information about an be found at [Metrics] and [ResourceUsage].

## Set up Cluster

``` bash
kwokctl create cluster --enable metrics-server
```

## Create Node and Pod

``` bash
kwokctl scale node --replicas 2
kwokctl scale pod --replicas 8
```

## Test Metrics

Wait about 45 seconds for the metrics server to collect the data.

``` bash
kubectl top node
kubectl top pod
```

[Metrics]: {{< relref "/docs/user/metrics-configuration" >}}
[ResourceUsage]: {{< relref "/docs/user/resource-usage-configuration" >}}
