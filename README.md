# **The Knative Approach to Serverless**

Repository dedicated to the practical work carried out for my master's degree thesis in **Engineering in computer science** at **Sapienza - Università di Roma**

Below you will find instructions for preparing a Kubernetes cluster using **K3S** as the chosen distribution and **Istio** as the ingress controller.

For installing Knative please refer to the guide in the [knative](knative) folder.

# Instructions

## Install k3s (disabling traefik)

    $ curl -sfL https://get.k3s.io | K3S_KUBECONFIG_MODE="644" sh -s - --disable traefik
***

## Install Istio

Install istioctl

    $ curl -sL https://istio.io/downloadIstioctl | sh -
    $ export PATH=$HOME/.istioctl/bin:$PATH

check before install Istio

    $ istioctl x precheck

if "Error: failed to get the Kubernetes version: Get "http://localhost:8080/version?timeout=5s": dial tcp 127.0.0.1:8080: connect: connection refused" appear you must copy

    $ sudo cp /etc/rancher/k3s/k3s.yaml $HOME/.kube/config

Enable auto-completion

- download full Istio release

        $ curl -L https://istio.io/downloadIstio | sh -
    
    then copy the "istio-x.xx.x/tools/istioctl.bash" file in home directory
- add these lines to ".bashrc" file
 
        export PATH=$HOME/.istioctl/bin:$PATH
        source ~/istioctl.bash

Install Istio with default parameters

    $ istioctl install -y

## install operator

Install with

    $ kubectl apply -f operatorEdited.yaml

check with

    $ kubectl get deployment knative-operator

***

## install knative-serving

check ingress with

    $ kubectl get svc istio-ingressgateway -n istio-system
    $ kubectl get pods -n istio-system

install knative-serving with DNS

    $ kubectl apply -f serving_personalDNS.yaml

then check with

    $ kubectl -n knative-serving get deployment
    $ kubectl -n knative-serving get KnativeServing knative-serving
***

## install knative-eventing

install with

    $ kubectl apply -f eventing.yaml

and check with

    $ kubectl get deployment -n knative-eventing
    $ kubectl get KnativeEventing knative-eventing -n knative-eventing

## configure autoTLS

install cert-manager

    $ kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.10.1/cert-manager.yaml

install cmctl

    $ OS=$(go env GOOS); ARCH=$(go env GOARCH); curl -fsSL -o cmctl.tar.gz https://github.com/cert-manager/cert-manager/releases/latest/download/cmctl-$OS-$ARCH.tar.gz
    $ tar xzf cmctl.tar.gz
    $ sudo mv cmctl /usr/local/bin

and check cert API with

    $ cmctl check api

create ClusterIssuer

    $ kubectl apply -f  cluster_issuer.yaml

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

turn only autoTLS, not both autoTLS and https redirect. There is a race-condition for the acme container that try to receive the ack in https with redirected active by default in the gateway

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

## NB: praticamente per colpa della race condition devi mettere redirected sul servizio solo dopo che il certificato è stato ottenuto, quindi la prima volta fai deploy senza http-protocol

and check with

    $ kubectl get configmap config-network -n knative-serving -o yaml

## deploy service with https redirected

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

## Remove the Knative Serving component

    $ kubectl delete KnativeServing knative-serving -n knative-serving

## Remove Knative Eventing component

    $ kubectl delete KnativeEventing knative-eventing -n knative-eventing

## Remove Operator

    $ kubectl delete -f operatorEdited.yaml