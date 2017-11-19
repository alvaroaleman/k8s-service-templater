# k8s-service-templater

## Description

A simple tool that can be used to render a template based on Kubernetes Services.

Included is a HAProxy template for the example usecase "For each NodePort service open an ipv6 socket and redirect it to its ipv4 counterpart".

## Usage

`k8s-service-templater -kubeconfig=kubeconfig -config=config.yml`
