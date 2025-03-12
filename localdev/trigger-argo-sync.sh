#!/bin/sh
argo_app_file="testapp.yaml"

# Update the testapp.yaml file
if [[ -f "$argo_app_file" ]]; then
    # Swap targetRevision between guestbook-v0.1 and guestbook-v0.2
    current_revision=$(yq eval ".spec.source.targetRevision" $argo_app_file)
    if [[ "$current_revision" == "HEAD" ]]; then
        new_revision="guestbook-v0.2" # Arbitrary tag in argocd guestbook repo that will cause indefinite syncing
    else
        new_revision="HEAD" # Move to revision that actually works
    fi

    yq eval "(.spec.source.targetRevision) |= \"$new_revision\"" -i testapp.yaml
    echo "Swapped targetRevision to $new_revision in $argo_app_file"
else
    echo "Error: $argo_app_file not found"
    exit 1
fi

kubectl apply -f $argo_app_file