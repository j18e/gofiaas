# Design

Since gofiaas attempts to immitate fiaas-deploy-daemon with all the newest
apiversions turned on, this document will attempt to enumerate the logic in FDD.

## deployer/kubernetes
This is the part of the Python app where most of the action happens. There is an
init script which sets up the dependency injection for all of:
- service
- service account
- ingress
- hpa
- ingress tls
- owner refs

### deployer/kubernetes/adapter.py
this will probably be the Deployer type in gofiaas
- sets up a K8s with all the deployers from the init script
- sets the version
- has a deploy method which calls deploy on everything based on an appspec
- has a similar delete method
- has a make_labels method which is used by deploy
  * app
  * fiaas/version
  * fiaas/deployment_id
  * fiaas/deployed_by
  * teams.fiaas labels
  * tags.fiaas labels

### deployer/kubernetes/deployment
Kubernetes deployments get their own folder since they're the most complex part
of deploying FIAAS applications. The init script sets up the following using
dependency injection:
- datadog
- prometheus
- generic_init_secrets
- deployment_secrets
- deployment_deployer

#### deployer/kubernetes/deployment/prometheus.py
If the given app spec has prometheus enabled, it adds annotations to the
deployment spec:
- prometheus.io/scrape
- prometheus.io/port
- prometheus.io/path

#### deployer/kubernetes/deployment/datadog.py
Does some setup in the init method that will apply to all sidecars created. The
apply method, if datadog is enabled in the given app spec:
- adds a datadog container to the deployment spec
- adds STATSD_HOST and STATSD_PORT env vars to the main container in the
  deployment

#### deployer/kubernetes/deployment/secrets.py
TODO detail out all this logic

#### deployer/kubernetes/deployment/deployer.py
Handles the actual managing of Kubernetes deployments.
