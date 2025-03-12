#!/bin/sh
# Please only run on locally installed K8s cluster for testing
POD_STATUS=$(kubectl get pods -n argocd -l app.kubernetes.io/name=argocd-server -o jsonpath='{.items[0].status.phase}')

if [[ "$POD_STATUS" != "Running" ]]; then
  kubectl create namespace argocd
  kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
  kubectl patch svc argocd-server -n argocd -p '{"spec": {"type": "LoadBalancer"}}'
  ARGO_PORT_1=$(kubectl get svc argocd-server -n argocd -o jsonpath='{.spec.ports[0].port}')
  ARGO_PORT_2=$(kubectl get svc argocd-server -n argocd -o jsonpath='{.spec.ports[1].port}')
  echo "ArgoCD installed, accessible on http://localhost:${ARGO_PORT_1} and https://localhost:${ARGO_PORT_2}"
else
  echo "ArgoCD is already running in the 'argocd' namespace"
  ARGO_PORT_1=$(kubectl get svc argocd-server -n argocd -o jsonpath='{.spec.ports[0].port}')
  ARGO_PORT_2=$(kubectl get svc argocd-server -n argocd -o jsonpath='{.spec.ports[1].port}')
  echo "Visit http://localhost:${ARGO_PORT_1} or https://localhost:${ARGO_PORT_2}"
fi