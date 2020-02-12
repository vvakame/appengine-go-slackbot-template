#!/bin/bash -eux

cd "$(dirname "$0")"

set +x

envsubst <app.yaml >app.generated.yaml

gcloud --quiet --project "$APPENGINE_APPLICATION" --verbosity warning app deploy ./app.generated.yaml
