apiVersion: v1
kind: Config
clusters:
- name: "bot-plugins"
  cluster:
    server: "https://rancher2.linecorp-dev.com/k8s/clusters/c-skhf4"
    api-version: v1

users:
- name: "user-r7kd7"
  user:
    token: "${USER_CREDENTIALS}"

contexts:
- name: "bot-plugins"
  context:
    user: "user-r7kd7"
    cluster: "bot-plugins"

current-context: "bot-plugins"
