# {{PROJECT}} — runbook

## Deploy

```sh
make docker-build
kubectl apply -k deployments/overlays/staging   # or: make k8s-deploy KUSTOMIZE_ENV=staging
kubectl apply -k deployments/overlays/prod      # or: make k8s-deploy KUSTOMIZE_ENV=prod
```

CD runs on push to `main` (`.github/workflows/cd.yml`); releases on `v*` tags
(`release.yml`).

## Rollback

```sh
kubectl rollout undo deployment/{{PROJECT}} -n {{PROJECT}}-prod
# or redeploy the previous image tag
```

## Monitoring & logs

- Health endpoint: `/healthz` (readiness/liveness probes in `deployments/deployment.yaml`)
- Logs: `kubectl logs -n {{PROJECT}}-prod deploy/{{PROJECT}} -f`

## Common issues

| Symptom | Likely cause | Fix |
|---------|--------------|-----|
| TODO | | |

## On-call

- TODO: escalation path, dashboards, alert channels
