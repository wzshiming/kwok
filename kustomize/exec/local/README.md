# Exec Local

This ClusterExec makes `kubectl exec` work for all Pod(s) managed by `kwok`.

There is no real container behind a fake Pod, so the command runs locally where `kwok` is running,
which is the host for `kwokctl` with the binary runtime, or the `kwok-controller` container otherwise.

Please refer to [Exec](https://kwok.sigs.k8s.io/docs/user/exec-configuration) for more on how it works.
