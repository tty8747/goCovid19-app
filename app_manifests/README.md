## Get env variables:
```bash
cd app_manifests/api_config/ && ./gen_config.sh && cd -
source env
```

## Create namespace:
```bash
kubectl create namespace "${NAMESPACE}" --dry-run=client -o yaml > namespace.yml
```

Test db connection:  
```bash
kubectl --namespace "${NAMESPACE}" run mysql --image mariadb:latest \
  --env="DB_ENTRYPOINT=${DB_ENTRYPOINT}" \
  --env="MARIADB_USER=${MARIADB_USER}" \
  --env="MARIADB_PASSWORD=${MARIADB_PASSWORD}" \
  --rm --stdin --tty -- bash
```
```bash
mysql --protocol=tcp "${MARIADB_DATABASE}" -u"${MARIADB_USER}" -p"${MARIADB_PASSWORD}" -h"${DB_ENTRYPOINT}"
```

## Create secret for goose:
```bash
kubectl create secret generic --namespace "${NAMESPACE}" goose-connection-string \
  --from-literal="GOOSE_DRIVER=mysql" \
  --from-literal="GOOSE_DBSTRING=${MARIADB_USER}:${MARIADB_PASSWORD}@tcp(${DB_ENTRYPOINT}:3306)/${MARIADB_DATABASE}" \
  --dry-run=client -o yaml > goose_secrets.yml
```

Tesing goose up command:
```bash
kubectl --namespace "${NAMESPACE}" run goose --image "${IMAGE_GOOSE}" \
       --overrides='
       {"spec": {
         "containers":
         [{
           "image": "tty8747/goose:latest",
           "name": "goose-latest",
           "command": ["goose"],
           "args": ["up"],
           "envFrom": [{
             "secretRef": {
               "name": "goose-connection-string"
               }
           }],
           "volumeMounts": [{
              "mountPath": "/home/store",
              "name": "store"
            }]
         }],
         "volumes":
         [{
           "name":"store",
           "emptyDir":{}
         }]
       }}' \
       --dry-run=client -o yaml
```

Check env inside goose container:
```bash
kubectl --namespace "${NAMESPACE}" run goose --image "${IMAGE_GOOSE}" \
       --overrides='
       {"spec": {
         "containers":
         [{
           "image": "tty8747/goose:latest",
           "name": "goose-latest",
           "command": ["/bin/sh", "-c", "--"],
           "args": ["while true; do sleep 30; done;"],
           "envFrom": [
             {"secretRef": {
                 "name": "goose-connection-string"
               }
             }
           ],
           "volumeMounts": [{
              "mountPath": "/home/store",
              "name": "store"
            }]
         }],
         "volumes":
         [{
           "name":"store",
           "emptyDir":{}
         }]
       }}'
```
Get env:
```bash
kubectl --namespace gocovid exec --tty --stdin goose -- env
```

Create config for api:
```bash
kubectl create configmap apiconfig --namespace "${NAMESPACE}" \
  --from-file=./api_config/config.yml \
  --dry-run=client -o yaml
```

## Create config with secrets for api:
```bash
kubectl create secret generic --namespace "${NAMESPACE}" api-secret-data \
  --from-file=./api_config/config.yml \
  --dry-run=client -o yaml > api_secret_data_config.yml
```

Create sensity data for api config:
```bash
kubectl create secret generic --namespace "${NAMESPACE}" api-connection-data \
  --from-literal="DB_ENTRYPOINT=${DB_ENTRYPOINT}" \
  --from-literal="MARIADB_DATABASE=${MARIADB_DATABASE}" \
  --from-literal="MARIADB_USER=${MARIADB_USER}" \
  --from-literal="MARIADB_PASSWORD=${MARIADB_PASSWORD}" \
  --dry-run=client -o yaml
```

## Create api deployment:
```bash
kubectl create deployment "${DEPLOYMENT}" --image "${IMAGE_API}" --port 5000 --replicas "${REPLICAS}" --namespace "${NAMESPACE}" --dry-run=client -o yaml > api_deploy.yml
```

### Also add into `api_deploy.yml` these uncomment lines:
```bash
# secrets into spec.containers
# envFrom:
# - secretRef:
#     name: api-connection-data

# config as volume into spec.containers for configMap
# volumeMounts:
# - name: api-config-volume
#   mountPath: /app/configs

# volume into spec.volumes for configMap
# - name: api-config-volume
#   configMap:
#     name: apiconfig

# config as volume into spec.containers for secrets
volumeMounts:
- name: api-secret-volume
  mountPath: /app/configs
  readOnly: true

# volume into spec.volumes for secrets
volumes:
- name: api-secret-volume
  secret:
    secretName: api-secret-data

# init containers into spec.initContainers
initContainers:
  - name: goose
    image: tty8747/goose:62bde2fe797e42878632774408c2e226a034aa0a
    command: ["goose"]
    args: ["up"]
    envFrom:
    - secretRef:
        name: goose-connection-string
# - image: tty8747/api:62bde2fe797e42878632774408c2e226a034aa0a
#   name: api-fill
#   command: ["/app/api"]
#   args: ["-fill=true"]
#   volumeMounts:
#   - name: api-secret-volume
#     mountPath: /app/configs
#     readOnly: true

# add into spec.containers.resources
resources:
  requests:
    memory: "64Mi"
    cpu: "250m"
  limits:
    memory: "128Mi"
    cpu: "900m"
```

## Expose api port through service:
```bash
kubectl expose deployment "${DEPLOYMENT}" --name "${SERVICE}" --port 5000 --target-port 5000 --protocol TCP --type=NodePort --namespace "${NAMESPACE}" --dry-run=client -o yaml > api_service.yml
```

## Create job for fill data"
```bash
kubectl --namespace gocovid create job "${JOB_NAME}" --image="${IMAGE_API}" --dry-run=client -o yaml
```

## Add into `spec.containers`:
```bash
# Add into spec.containers
command: ["/app/api"]
args: ["-fill=true"]
volumeMounts:
- name: api-secret-volume
  mountPath: /app/configs
  readOnly: true

# volume into spec.volumes for secrets
volumes:
- name: api-secret-volume
  secret:
    secretName: api-secret-data

# add into spec 
restartPolicy: OnFailure
```

## Create frontend deployment:
```bash
kubectl create deployment "${FRONT_DEPLOYMENT}" --image "${IMAGE_FRONT}" --port 8080 --replicas "${REPLICAS}" --namespace "${NAMESPACE}" --dry-run=client -o yaml > front_deploy.yml
```

### Also add into `front_deploy.yml` these uncomment lines:
```bash
# add into spec.containers
command: ["/app/web"]
args: [ "-addr", "0.0.0.0:8080", "-apiVers", "v1", "-hostname", "svc-gocovid", "-port", "5000"]

# add into spec.containers.resources
resources:
  requests:
    memory: "64Mi"
    cpu: "250m"
  limits:
    memory: "128Mi"
    cpu: "900m"
```

## Expose front port through service:
```bash
kubectl expose deployment "${FRONT_DEPLOYMENT}" --name "${FRONT_SERVICE}" --port 8080 --target-port 8080 --protocol TCP --type=NodePort --namespace "${NAMESPACE}" --dry-run=client -o yaml > front_service.yml > job_fill.yml
```

## Create ingress:
```bash
# set your LOADBALANCER_NAME variable
kubectl create ingress "${INGRESS}" --namespace="${NAMESPACE}" --annotation kubernetes.io/ingress.class=alb --annotation alb.ingress.kubernetes.io/scheme=internet-facing --annotation alb.ingress.kubernetes.io/load-balancer-name="${LOADBALANCER_NAME}" --annotation alb.ingress.kubernetes.io/target-group-attributes=stickiness.enabled=true --annotation stickiness.lb_cookie.duration_seconds=60 --annotation alb.ingress.kubernetes.io/target-type=ip --rule="/*=${FRONT_SERVICE}:8080" --dry-run=client -o yaml > ingress.yml
```

## Create horizontal pod autoscaller:
```bash
kubectl autoscale --namespace "${NAMESPACE}" deployment "${DEPLOYMENT}" --cpu-percent=70 --min=1 --max=50 --dry-run=client -o yaml > api_hpa.yml
kubectl autoscale --namespace "${NAMESPACE}" deployment "${FRONT_DEPLOYMENT}" --cpu-percent=70 --min=1 --max=50 --dry-run=client -o yaml > front_hpa.yml
```

## Add in both hpa into `HorizontalPodAutoscaler.metadata` section:
```bash
namespace: gocovid
```

## Test hpa:
```bash
docker run --rm --net=host loadimpact/loadgentest-wrk -c 100 -t 100 -d 15m http://test.app.ubukubu.ru
```

```bash
# kubectl --namespace gocovid logs --follow deploy-gocovid-57f4d5ddf6-6hxpd -c goose
```
