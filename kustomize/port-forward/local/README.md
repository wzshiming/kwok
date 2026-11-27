# Port Forward Local

This ClusterPortForward makes `kubectl port-forward` work for all Pod(s) managed by `kwok`.

There is no real container behind a fake Pod, so every port of the Pod is forwarded to the server of `kwok` itself,
which listens on `127.0.0.1:10247` by default and serves endpoints such as `/healthz` and `/metrics`.

Please refer to [Port Forward](https://kwok.sigs.k8s.io/docs/user/port-forward-configuration) for more on how it works.
