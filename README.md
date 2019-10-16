## Welcome to "hello" application!
### some choice
`python` or `golang`  
I like to do tasks in Python. But here we require production solution and I think Python could add some configuration complexity with WSGI/FastCGI. Golang is a single binary option and seems easy to deploy. So here we `go` :)  

Kubernetes in GCP or AWS  
AWS ECS  
AWS EC2  

`AWS` or `GCP`  

`storage/database` choice  
Simple interface with http json API  

assumptions
`<username>` is unique  

NB For `high load` application mode cluster will scale via Autoscaling Group. AMI image I would create with Packer. `etcd` supports service discovery.  
NB For `secure` setup I would use authentication and harden endpoint with TLS.  

Sometimes it is not an easy decision to make while designing infrastructure.  
Because it will affect further maintenance and operational costs, velocity of development cycle.  
It is not a single person decision.  
I would like to ask you guys to discuss together pros and cons, and make a deliberate decision ;)  

### application
i took liberties with requirements and added two extra endpoints  
`/hc` because it is needed by loadbalancer infrastructure, and  
`/version` this one helps to debug and test deployment rollout  

### infrastructure
export the following environment variables and set parameters in `var.tf` file.  
environment variables consist of sevsitive data, while var.tf consists of parameters committed to the source.  
```
export AWS_ACCESS_KEY_ID=<your aws access key>
export AWS_SECRET_ACCESS_KEY=<your aws secret key>
export AWS_DEFAULT_REGION=<aws region of operation>
export TF_VAR_region=$AWS_DEFAULT_REGION
# the latter duplication is necessary for work together terraform and aws cli
# aws cli needed later during build stage
```
in order to deploy cloud infrastructure to aws we use terraform tool and HCL configuration language.  
terraform v0.12 needs to be installed on management workstation.  
```
curl -LO https://releases.hashicorp.com/terraform/0.12.10/terraform_0.12.10_linux_amd64.zip
unzip ./terraform_0.12.10_linux_amd64.zip
sudo mv ./terraform /usr/local/bin
```
checkout master branch, initialize your working directory and deploy infrastructure.  
NB in terraform we use `local` backend. this works for single user mode. to work on infrastructure collaboratively we require to setup remote (shared) backend. before we used s3, but hashicorp provides cloud backend now. no needs to provision buckets in advance anymore.
```
git checkout master
terraform init
terraform apply
```

### build
install aws cli
```
pip3 install awscli --upgrade --user
```
export the following environment variable either in local or github actions environment.  
```
terraform state show aws_ecr_repository.ecr | grep repository_url
export HELLO_REPOSITORY_URL=<repository_url>
cd src
./build.sh
```
github actions
ECR

### deployment
Intentionally manual process. Based on my previous experience all production releases were semi-automated.
Because of existence of release management process.
Goal here is to automate deployment, having a single action.
In order to deploy `hello` to production we need two actions `bump app version` & `terraform apply`.
```
# this step reads the latest artifact version, and set it inside infra configuration
# for the sake of simplicity we use commit SHA for artifacts versioning
./bump.sh

# release the latest app version
terraform apply
# plan should say the following
#         Plan: 1 to add, 1 to change, 1 to destroy.
# it is correct, we are replacing ecs task, and updating in-place service. type "yes"
```

### integration test
```
terraform state show aws_lb.alb | grep dns_name
curl <dns_name>:5000/version | jq
```

### local deployment and tests
```
# install golang
curl -O https://dl.google.com/go/go1.13.1.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.13.1.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
# install compose
sudo curl -L "https://github.com/docker/compose/releases/download/1.24.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
# run redis
docker-compose up -d
# unit test
cd src
go test -v
# build app
go build -o hello
# run app
export REDIS_ENDPOINT=localhost:6379
./hello 5000 local
# integration test
curl localhost:5000/version | jq
```
