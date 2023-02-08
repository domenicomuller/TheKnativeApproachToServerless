# **The Knative Approach to Serverless**

Repository dedicated to the practical work carried out for my master's degree thesis in **Engineering in computer science** at **Sapienza - Università di Roma**

Below you will find instructions for preparing a Kubernetes node using **K3S** as the chosen distribution and **Istio** as the ingress controller.

For Knative installation please refer to the guide in the [knative](knative) folder.

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

# Uninstall

If you installed K3s using the installation script, a script to uninstall K3s was generated during installation. To uninstall K3s from a server node, run:

    $ /usr/local/bin/k3s-uninstall.sh