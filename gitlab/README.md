# Install gitlab

## First off, we need to create an infrastructure by Terraform

- `terraform init`
- `terraform plan`
- `terraform apply --auto-approve`

## In the second place, we need to set up our infrastructure

- Check the host is alive:  
  `ansible -i ansible/inventory/hosts.yml all -m raw -a "cat /etc/*release"`
- Install requirements:  
  `ansible-galaxy install -r ansible/requirements.yml`
- At the first run, `remove_gitlab_data` variable should set true:  
  `ansible-playbook -i ansible/inventory/hosts.yml ansible/main.yml -e "remove_gitlab_data=true"`
- Get root pass:  
  `ansible-playbook -i ansible/inventory/hosts.yml ansible/main.yml --tags "get_rootpass"`

### Scenario when we lost GitLab node in some region

- Lose our node:  
  `terraform destroy -target aws_instance.gitlab`
- Create new node in another region:  
  `terraform plan -var="subnet_zone=eu-central-1c"`  
  `terraform apply -var="subnet_zone=eu-central-1c"`
- Ensure that `remove_gitlab_data` variable is false if you want to save gitlab data:  
  `ansible-playbook -i ansible/inventory/hosts.yml ansible/main.yml`

## Runner

### Start runner
```bash
docker run -d --name gitlab-runner --restart always \
     -v /srv/gitlab-runner/config:/etc/gitlab-runner \
     -v /var/run/docker.sock:/var/run/docker.sock \
     gitlab/gitlab-runner:latest
```

### Register runner
```bash
GITLAB_REGISTRATION_TOKEN="xxxxxxxxxxxxxx65oxH"

# The Docker untagged executer
docker run --rm -it -v /srv/gitlab-runner/config:/etc/gitlab-runner gitlab/gitlab-runner register -n --url https://gitlab.ubukubu.ru/ --registration-token "${GITLAB_REGISTRATION_TOKEN}" --executor docker --description "aws runner" --docker-image ubuntu:latest --run-untagged

# The Docker executor with some tags
RUNNER_TAG_LIST="golang, go"
docker run --rm -it -v /srv/gitlab-runner/config:/etc/gitlab-runner gitlab/gitlab-runner register -n --url https://gitlab.ubukubu.ru/ --registration-token "${GITLAB_REGISTRATION_TOKEN}" --executor docker --description "hetzner runner" --tag-list "${RUNNER_TAG_LIST}" --docker-image ubuntu:latest

# The Docker executor with Docker socket binding
# https://docs.gitlab.com/ee/ci/docker/using_docker_build.html#use-the-docker-executor-with-docker-socket-binding
RUNNER_TAG_LIST="docker.sock"

docker run --rm -it -v /srv/gitlab-runner/config:/etc/gitlab-runner \
     gitlab/gitlab-runner register -n --url https://gitlab.ubukubu.ru/ \
     --registration-token "${GITLAB_REGISTRATION_TOKEN}" --executor docker \
     --docker-image "docker:19.03.12" --tag-list "${RUNNER_TAG_LIST}" \
     --description "docker sock runner" \
     --docker-volumes /var/run/docker.sock:/var/run/docker.sock
```

### List runners
```bash
docker run --rm -it -v /srv/gitlab-runner/config:/etc/gitlab-runner gitlab/gitlab-runner list
```

### Unregister runners
```bash
docker run --rm -it -v /srv/gitlab-runner/config:/etc/gitlab-runner gitlab/gitlab-runner unregister --all-runners
```
