# Get k8s cluster credential
gcloud container clusters get-credentials pilot-cluster

# Connect to mysql
# Start this proxy in one terminal
cloud_sql_proxy -instances=pilot-154919:asia-southeast1:url-shortener=tcp:33061 -credential_file=/Users/woraphol/pilot-154919-ca5f429d43f7.json &

# create a secret
kubectl create secret generic mysql --from-literal=mysql-username=url_shortener --from-literal=mysql-password=EZq65QhXnybDCea

# apply changes
kubectl apply -f ./k8s/url-shortener.yml

# Build project using cloudbuild
gcloud builds submit --config cloudbuild.yml . --substitutions=SHORT_SHA="1.7"

