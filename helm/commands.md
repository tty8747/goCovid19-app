```bash
# https://docs.gitlab.com/ee/user/packages/helm_repository/

helm template gocovid --set dbSettings.name=dbdbdb
helm lint gocovid --set dbSettings.name=dbdbdb
helm package gocovid
helm package gocovid --app-version 2.2.3 --version 0.2.0

# Get deploy token: Settings -> Repository -> Deploy Tokens
TOKEN_USER="gitlab+deploy-token-978533"
TOKEN_PASS="VV1Evaxxxxxxxxxxxxxx"
CHANNEL="stable"
PROJECT_ID="34764615"
CI_API_V4_URL="https://gitlab.ubukubu.ru/api/v4"

curl --request POST \
     --form 'chart=@gocovid-0.1.0.tgz' \
     --user "${TOKEN_USER}:${TOKEN_PASS}" \
     "${CI_API_V4_URL}/projects/${PROJECT_ID}/packages/helm/api/${CHANNEL}/charts"

curl --request POST \
     --form 'chart=@gocovid-0.2.0.tgz' \
     --user "${TOKEN_USER}:${TOKEN_PASS}" \
     "${CI_API_V4_URL}/projects/${PROJECT_ID}/packages/helm/api/${CHANNEL}/charts

helm repo add --username "${TOKEN_USER}" --password "${TOKEN_PASS}" gocovid-repo "${CI_API_V4_URL}/projects/${PROJECT_ID}/packages/helm/${CHANNEL}"

# Put your real data:
helm install gocovid gocovid-repo/gocovid \
  --set front.albName=alb-eks-test-sg \
  --set dbSettings.endpoint=gocovid-test-2022042309290953210000000b.c9nmurf5weua.eu-central-1.rds.amazonaws.com \
  --set dbSettings.name=gocovid \
  --set dbSettings.user=testuser \
  --set dbSettings.password=yourstrongpassword

helm upgrade gocovid gocovid-repo/gocovid \
  --set front.albName=alb-eks-test-sg \
  --set dbSettings.endpoint=gocovid-test-2022042309290953210000000b.c9nmurf5weua.eu-central-1.rds.amazonaws.com \
  --set dbSettings.name=gocovid \
  --set dbSettings.user=testuser \
  --set dbSettings.password=yourstrongpassword

helm history gocovid
helm rollback --dry-run gocovid 1
helm rollback  gocovid 1
```
