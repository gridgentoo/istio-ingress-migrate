# istio-ingress-migrate

This tool helps users migrate from older versions of Istio (1.5 and older) to newer versions when using the Kubernetes Ingress resource.

Prior to Istio 1.6, HTTPS usage for Ingress required configuring both an Istio `Gateway` and a Kubernetes `Ingress`; the `tls` section of the `Ingress` spec was ignored.

With Istio 1.6+, the [full Ingress spec is natively supported](https://istio.io/latest/docs/tasks/traffic-management/ingress/kubernetes-ingress/.

This tool helps users migrate from a combination of `Gateway` and `Ingress` to just using `Ingress`.
An alternative option is to migrate to using `Gateway` and `VirtualService` which is documented on the [istio.io site](https://istio.io/latest/docs/tasks/traffic-management/ingress/ingress-control/).

## Usage

This tool will take `Ingress` and `Gateway` YAML and output converted `Ingress` YAML.

```shell
kubectl get ingresses.extensions,gateways.v1alpha3.networking.istio.io -A -oyaml | istio-ingress-migrate > ingress.yaml
```

This translation is not guaranteed to be possible; there are some features in `Gateway` not supported by `Ingress`.
In these cases, the command will fail with a non-zero exit code and output warnings indicating what configurations need to be manually resolved.

Once converted, the new Ingress's can be applied to the cluster. This can be done before upgrading to Istio 1.6, as the `tls` field is ignored in the new version.
Upon upgrade, the old `Gateway`s can be ignored.
