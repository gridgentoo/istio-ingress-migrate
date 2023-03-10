Original repository.     
https://github.com/istio-ecosystem/istio-ingress-migrate.    

# istio-ingress-migrate

This tool helps users migrate from older versions of Istio (1.5 and older) to newer versions when using the Kubernetes Ingress resource.

Prior to Istio 1.6, HTTPS usage for Ingress required configuring both an Istio `Gateway` and a Kubernetes `Ingress`; the `tls` section of the `Ingress` spec was ignored.

With Istio 1.6+, the [full Ingress spec is natively supported](https://istio.io/latest/docs/tasks/traffic-management/ingress/kubernetes-ingress/.

This tool helps users migrate from a combination of `Gateway` and `Ingress` to just using `Ingress`.
An alternative option is to migrate to using `Gateway` and `VirtualService` which is documented on the [istio.io site](https://istio.io/latest/docs/tasks/traffic-management/ingress/ingress-control/).

## Installation

### Precompiled binaries

Precompiled binaries can be found in the [releases](https://github.com/istio-ecosystem/istio-ingress-migrate/releases) page.

### Using go

The latest version can be compiled directly with:

```
go install github.com/istio-ecosystem/istio-ingress-migrate@v0.0.2
```

## Usage

This tool will take `Ingress` and `Gateway` YAML and output converted `Ingress` YAML.

```shell
kubectl get ingresses.extensions,gateways.v1alpha3.networking.istio.io -A -oyaml | istio-ingress-migrate > ingress.yaml
```

This translation is not guaranteed to be possible; there are some features in `Gateway` not supported by `Ingress`.
In these cases, the command will fail with a non-zero exit code and output warnings indicating what configurations need to be manually resolved.

Once converted, the old Gateway's should be annotated to indicate the migration is complete:

```shell
kubectl annotate gateways.networking.istio.io istio-autogenerated-k8s-ingress -n gke-system gateway.istio-ecosystem.io/migrated=true
```

Once converted, the new Ingress's can be applied to the cluster. This can be done before upgrading to Istio 1.6, as the `tls` field is ignored in the new version.
Upon upgrade, the old `Gateway`s can be ignored.

## Error Dictionary

While the conversion logic attempts to automatically convert all resources, some conversions result in ambiguous
scenarios which cannot be resolved automatically.
When this occurs, an error will be logged to `stderr`, and the command will exit with a non-zero exit code.

Below list the possible errors and how to resolve them.

### ERR01: Unsupported Port

Ingress only allows port 80 (HTTP) and port 443 (HTTPS).
This error will be emitted when a Gateway is configured for another port.
For example:

```yaml
servers:
- hosts:
  - '*'
  port:
    name: http
    number: 1234
    protocol: HTTP
```

If you require custom ports, please see [Migrating to Istio Gateway](#migrating-to-istio-gateway).

If you do not require custom ports, this error can be ignored; the port will no longer be usable.

### ERR02: Missing TLS configuration

This error occurs when port 443 is defined on the Gateway, but TLS settings are not configured.
For example:

```yaml
servers:
- hosts:
  - '*'
  port:
    name: http
    number: 443
    protocol: HTTP
```

This configuration is not supported; port 443 can only be HTTPS.


If you require HTTP on port 443, please see [Migrating to Istio Gateway](#migrating-to-istio-gateway).

If you do not require HTTP on port 443, remove the offending configuration from the Gateway.

### ERR03: UnsupportedTLSSettings

Gateway supports some advanced TLS customizations, such as mutual TLS.
Ingress supports only specifying the certificate to use.
For example:

```yaml
servers:
- hosts:
  - mtls.example.com
  port:
    name: https
    number: 443
    protocol: HTTPS
  tls:
    credentialName: mtls-cred
    mode: MUTUAL
```

This configuration is not supported by Ingress.


If you require custom TLS settings, please see [Migrating to Istio Gateway](#migrating-to-istio-gateway).

If you do not require custom TLS settings, the error can be safely ignored; the settings will no longer apply after migration.

### ERR04: Host conflicts

This error occurs when there are conflicting settings defined in the Gateway. 
For example:

```yaml
servers:
- hosts:
  - 'service.com'
  tls:
    mode: SIMPLE
    credentialName: secret-a
- hosts:
  - 'service.com'
  tls:
    mode: SIMPLE
    credentialName: secret-b
```

One of the two configurations should be removed, then the tool can be re-run.

### ERR05: Existing TLS mismatch

Prior to migration, Ingress could define TLS settings, but they would be ignored.
This error occurs when the already defined TLS settings conflict with what is expected from the Gateway.

For example:

Ingress:
```yaml
spec:
  tls:
  - hosts: ["example.com"]
    secretName: my-ingress-credential
```

Gateway:
```yaml
servers:
- hosts:
  - example.com
  tls:
    credentialName: my-gateway-credential
    mode: SIMPLE
```

In this case, the same hostname is configured to use `my-ingress-credential` and `my-gateway-credential`.
Depending on which you expect to actually use, update the two fields to be refer to the same Secret.

Note: when using Istio 1.4, the Gateway's `credentialName` would be the one used, so in most cases the Ingress should be updated to match.

### ERR06: No credential found

This error occurs when the `host` in Ingress has no matching HTTPS configuration in a Gateway.
In this case, TLS will not be configured for the Ingress, meaning only HTTP requests will be allowed.

If you want HTTPS to be used for this Ingress, fill in the `tls` section in the Ingress or add configuration to the Gateway for this host.

If you do not need HTTPS for this Ingress, the error can be ignored; only HTTP will be allowed.

## Migrating to Istio Gateway

This tool focuses on migrating from Istio Gateway and Kubernetes Ingress to just Kubernetes Ingress.
Alternatively, you can do the opposite and migrate to using Istio Gateway and VirtualService.
The configuration to do this is out of scope for this tool.

To perform this migration, you can use `istioctl convert-ingress` (note: this command was removed in istioctl 1.9) and follow the [Istio documentation](https://istio.io/latest/docs/tasks/traffic-management/ingress/ingress-control/).
