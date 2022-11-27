# Ref: https://cloud.google.com/kubernetes-engine/docs/how-to/creating-an-autopilot-cluster

CLUSTER_NAME=aerospike-dev
REGION=us-west1
PROJECT_ID=aerospike-dev-367320

auth:
	gcloud auth login

build:
	gcloud container clusters create-auto "${CLUSTER_NAME}" \
    	--region "${REGION}" \
    	--project "${PROJECT_ID}";

connect:
	gcloud container clusters get-credentials "${CLUSTER_NAME}" \
    	--region "${REGION}" \
    	--project "${PROJECT_ID}";

delete:
	gcloud container clusters delete "${CLUSTER_NAME}" \
    	--region "${REGION}" \
    	--project "${PROJECT_ID}";
