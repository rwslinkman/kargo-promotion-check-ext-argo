# Kargo Promotion - check external Argo
KPCEA - Basic validation client to run Kargo `AnalysisTemplates` for verifying external Freight.  

## About
[Kargo](https://docs.kargo.io/user-guide/core-concepts/) is an interesting project that helps to promote your GitOps architecture from one environment to the other.  
It is made by the developers of [ArgoCD](https://argo-cd.readthedocs.io/en/stable/) and the applications work neatly together.  

Kargo bases its functionality on the idea that your deployable applications are living in the same cluster and are separated by namespace.  
This allows ArgoCD to manage all environments of your application.  
In [my experience](https://www.linkedin.com/in/rwslinkman/) as a Cloud/Software Engineer I find that companies will have an entirely separate cluster for their Development, Acceptance or UAT environment.  
Even though it is possible to allow one ArgoCD to manage resources in multiple clusters, I also find situations where each environment/cluster has its own ArgoCD instance.  

To make Kargo work in these cases, it will have to use the other (external) ArgoCD instances to verify the successful sync of your config to the next environment.  
This `go` application will perform said verifications on the provided ArgoCD server using the Argo API.  
Since the container will exit with statuscode 0 or 1, it can be used as a Kubernetes [Job](https://kubernetes.io/docs/concepts/workloads/controllers/job/).    
The `Job` in turn can be used in the `AnalysisTemplate` of Kargo.  

## Usage

### General
Create the following [AnalysisTemplate](https://docs.kargo.io/user-guide/how-to-guides/working-with-stages#verification) for Kargo.  
Please note that most config variables are hidden from this example. See table below.   
It is not recommended to use the `latest` version. Always use tagged versions.  
See [Docker Hub](https://hub.docker.com/r/rwslinkman/kargo-promotion-check-ext-argo/tags) for recent tags.  

```yaml
apiVersion: argoproj.io/v1alpha1
kind: AnalysisTemplate
metadata:
  name: check-external-argo
  namespace: namespace-of-kargo-app
spec:
  metrics:
  - name: check-external-argo
    provider:
      job:
        spec:
          template:
            spec:
              containers:
              - name: kpcea
                image: rwslinkman/kargo-promotion-check-ext-argo:latest
                env:
                  - name: ARGOCD_SERVER
                    value: argocd.mydomain.xyz
              restartPolicy: Never
          backoffLimit: 1
```

The container must be configured with a few parameters and has some optional config.  
These need to be set as environment variables.   

| Variable                | Description                        | Required    | Note                                                  |
|-------------------------|------------------------------------|-------------|-------------------------------------------------------|
| `ARGOCD_SERVER`         | The Argo CD server address         | Yes         | Remove protocol from URL when providing (no https://) |
| `ARGOCD_API_TOKEN`      | API token for authentication       | Conditional | Only required in TOKEN mode                           |
| `ARGOCD_APP_NAME`       | The Argo CD application name       | Yes         | n/a                                                   |
| `ARGOCD_API_USERNAME`   | Username for Argo CD API access    | Conditional | Only required in LOGIN mode                           |
| `ARGOCD_API_PASSWORD`   | Password for Argo CD API access    | Conditional | Only required in LOGIN mode                           |
| `KPCEA_TARGET_REVISION` | Target Git revision for deployment | Yes         | n/a                                                   |
| `KPCEA_TIMEOUT`         | Timeout duration (in seconds)      | No          | Defaults to `30` seconds                              |
| `KPCEA_INTERVAL`        | Sync interval (in seconds)         | No          | Defaults to `5` seconds                               |
| `KPCEA_INSECURE`        | Allow insecure connections         | No          | Defaults to `false`                                   |

### TOKEN mode vs. LOGIN Mode
KPCEA relies on a [local user from ArgoCD](https://argo-cd.readthedocs.io/en/stable/operator-manual/user-management/#create-new-user) to get access to the desired ArgoCD instance.  
This can be set directly in the `ARGOCD_API_TOKEN` parameter.   
If you are not able to get a token from ArgoCD, it is possible to use LOGIN mode on KPCEA.  

Provide the `ARGOCD_API_USERNAME` and `ARGOCD_API_PASSWORD` parameters and leave the `ARGOCD_API_TOKEN` empty.  
The KPCEA client will pick this up and retrieves a (temporary) token from ArgoCD for 1 session.   

