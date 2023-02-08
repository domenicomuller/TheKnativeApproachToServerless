# **Install Knative**

Below I will show the steps on how to install Knative on a ready Kubernetes node. The procedure will be followed through the use of the Knative Operator.

# Instructions

## Install operator

Install with

    $ kubectl apply -f https://github.com/knative/operator/releases/download/knative-v1.8.0/operator.yaml

check with

    $ kubectl get deployment knative-operator

***

## Install knative-serving

Check ingress with

    $ kubectl get svc istio-ingressgateway -n istio-system
    $ kubectl get pods -n istio-system

Install knative-serving with DNS

    $ kubectl apply -f serving_personalDNS.yaml

then check with

    $ kubectl -n knative-serving get deployment
    $ kubectl -n knative-serving get KnativeServing knative-serving
***

## Install knative-eventing

Install with

    $ kubectl apply -f eventing.yaml

and check with

    $ kubectl get deployment -n knative-eventing
    $ kubectl get KnativeEventing knative-eventing -n knative-eventing
***

## Configure autoTLS with cert-manager

install cert-manager

    $ kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.10.1/cert-manager.yaml

install cmctl

    $ OS=$(go env GOOS); ARCH=$(go env GOARCH); curl -fsSL -o cmctl.tar.gz https://github.com/cert-manager/cert-manager/releases/latest/download/cmctl-$OS-$ARCH.tar.gz
    $ tar xzf cmctl.tar.gz
    $ sudo mv cmctl /usr/local/bin

and check cert API with

    $ cmctl check api

create ClusterIssuer

    $ kubectl apply -f cluster_issuer.yaml

check ClusterIssuer

    $ kubectl get clusterissuer letsencrypt-issuer -o yaml

install cert-manager

    $ kubectl apply -f https://github.com/knative/net-certmanager/releases/download/knative-v1.8.0/release.yaml

check cert-manager

    $ kubectl get deployment net-certmanager-controller -n knative-serving

append in config-certmanager, under the "data" field, the ClusterIssuer ref

    $ KUBE_EDITOR="nano" kubectl edit configmap config-certmanager -n knative-serving

like

    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: config-certmanager
      namespace: knative-serving
      labels:
        networking.knative.dev/certificate-provider: cert-manager
    data:
      issuerRef: |
        kind: ClusterIssuer
        name: letsencrypt-issuer

check config-certmanager

    $ kubectl get configmap config-certmanager -n knative-serving -o yaml

turn only autoTLS, not both autoTLS and https redirect. There is a race-condition for the acme container that try to receive the ack in https with redirected active by default in the gateway. To avoid this, it is advisable to activate this property only after the correct generation of the certificate.

    $ KUBE_EDITOR="nano" kubectl edit configmap config-network -n knative-serving

like

    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: config-network
      namespace: knative-serving
    data:
       ...
       auto-tls: Enabled
       ...

and check with

    $ kubectl get configmap config-network -n knative-serving -o yaml

## Deploy service with https redirected

    apiVersion: serving.knative.dev/v1
    kind: Service
    metadata:
      name: helloworld-go-serving
      namespace: default
      annotations:
        # networking.knative.dev/disable-auto-tls: "true"
        networking.knative.dev/http-protocol: "redirected"
    spec:
      template:
        spec:
          containers:
          - image: docker.io/domll/helloworld-go-serving
            env:
            - name: TARGET
            value: "Go Sample v1"

# Uninstall

## Remove the knative-serving component

    $ kubectl delete KnativeServing knative-serving -n knative-serving

## Remove knative-eventing component

    $ kubectl delete KnativeEventing knative-eventing -n knative-eventing

## Remove Operator

    $ kubectl delete -f https://github.com/knative/operator/releases/download/knative-v1.8.0/operator.yaml