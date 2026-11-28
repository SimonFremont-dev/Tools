# Deployments

Kustomize manifests for {{PROJECT}}.

- `kustomization.yaml` + the base manifests (namespace, configmap, deployment,
  service, ingress, hpa)
- `overlays/dev|staging|prod/` — per-environment overlays (namespace, name
  prefix, replica counts, image tag)

Deploy:

```sh
kubectl apply -k deployments/overlays/dev     # or staging / prod
make k8s-deploy KUSTOMIZE_ENV=prod
```

Secrets are never committed: use Sealed Secrets or External Secrets and
reference them from the deployment (see the `envFrom`/`env` sections).
