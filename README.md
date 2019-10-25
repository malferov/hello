## Welcome to "hello" application!
### requirements
Create a sample application with http json API, production ready infrastructure, and automated one click deployment.  
Diagram  
https://docs.google.com/document/d/1AdzaqEfxJR_i9JV5xn-o4WkYh_mcCfjjXyZN1PIPRFc/edit?usp=sharing

### some choice
`python` or `golang`  
I like to do tasks in Python. But here we require production solution and I think Python could add some configuration complexity with WSGI/FastCGI. Golang is a single binary option and seems easy to deploy. So here we `go` :)  

`kubernetes or not`  
I was thinking of implementing in Kubernetes. A bit complex to setup, but very handy deployments. I did not choose this option because I don't have enough production experience with Kube.  

`AWS` or `GCP`  
I would choose GCP for kubernetes setup definitely. In AWS EKS we have to manage cluster nodes ourselves. But for docker service AWS is more commonly used. So ECS.  

`storage/database` choice  
I was looking for simple interface with http json API. The purpos was to have simple app without extra SDK dependency. I was considering CouchDB and Elasticsearch. Later is top product but not exactly that we want. We need key/value store. CouchDB matches well, but needs to be installed on EC2 (less automation), or container with persistent volume (more configs). I also was considering Etcd and KV engine for Consul or Vault. All those have sophisticated cluster setup. I did not like all options above, and realized we have in AWS Redis managed service. That option works for both cloud and on premises. And sdk is quite light for go. So here we go.  

### application
I had previously experience with gin-gonic framework. Some time ago I chose that framework because of its excellent performance reports.  
The following two extra endpoints will be added to application for utility purpose  
`/hc` because it is needed by load balancer infrastructure, and  
`/version` this one helps to debug and test deployment rollout  

### infrastructure
I was bearing in mind the following aspects while creating infrastructure: `security`, `reliability`, `performance`, `operational efficiency`, `cost aspect`.  
`NB` For `high load` application mode, autoscaling scheduler needs to be implemented. This setup is fixed in size. Amount of ECS and Redis nodes, as well as application instances could be adjusted in `var.tf`  
`NB` For `secure` setup I would use authentication and harden endpoints with TLS.  

Use root user in your AWS account, or manually create `management_user` and apply the policy below.
<details><summary>IAM policy for management user</summary>

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "iam:*",
            "Resource": [
                "arn:aws:iam::*:policy/*",
                "arn:aws:iam::*:user/*",
                "arn:aws:iam::*:group/*",
                "arn:aws:iam::*:role/*",
                "arn:aws:iam::*:instance-profile/*"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "ec2:*"
            ],
            "Resource": "*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "autoscaling:*",
                "elasticloadbalancing:*"
            ],
            "Resource": "*"
        },
        {
            "Effect": "Allow",
            "Action": "ecs:*",
            "Resource": "*"
        },
        {
            "Effect": "Allow",
            "Action": "ecr:*",
            "Resource": "*"
        },
        {
            "Effect": "Allow",
            "Action": "elasticache:*",
            "Resource": "*"
        }
    ]
}
```

</details>

After creation management user and applying policy above, create security credentials for that user.  
In order to deploy infrastructure, proceed with following steps on Management Workstation.  
Export the following environment variables and set parameters in `var.tf` file.  
Environment variables consist of sensitive data, while var.tf consists of parameters committed to the source.  
Please note it is necessary to whitelist your nat public ip in var.tf (`trusted` variable)
```
export AWS_ACCESS_KEY_ID=<your_aws_access_key>
export AWS_SECRET_ACCESS_KEY=<your_aws_secret_key>
export AWS_DEFAULT_REGION=<aws_region_of_operation>
export TF_VAR_region=$AWS_DEFAULT_REGION
# the latter duplication is necessary for work together terraform and aws cli
# aws cli needed later during build stage
```
In order to deploy cloud infrastructure to aws we use terraform tool and HCL configuration language.  
terraform v0.12 needs to be installed on Management Workstation.  
```
curl -LO https://releases.hashicorp.com/terraform/0.12.10/terraform_0.12.10_linux_amd64.zip
unzip ./terraform_0.12.10_linux_amd64.zip
sudo mv ./terraform /usr/local/bin
```
checkout master branch, initialize your working directory and deploy infrastructure.  
`NB` In terraform we use `local` backend. This works for single user mode. In order to work on infrastructure collaboratively we require to setup remote (shared) backend. Earlier we used s3, but hashicorp provides cloud backend now. No needs to provision buckets in advance anymore.
```
git checkout master
terraform init
terraform apply
```

### build
There is no branch model, unit tests run only on master. It is triggered by GitHub Actions. Whenever we want to build an image from any tag or commit, we pull that version to Management Workstation.
Then build an artifact and push in ECR. Later the artifact is used by service configuration in ECS.  
On Management Workstation  
```
# install aws cli
pip3 install awscli --upgrade --user
```
export the following environment variable  
```
terraform state show aws_ecr_repository.ecr | grep repository_url
export HELLO_REPOSITORY_URL=<repository_url>
```
build application and push artifact  
```
cd src
./build.sh
```

### app deployment
Intentionally manual process. Based on my previous experience all production releases were semi-automated.
Because of existence of release management process.
Goal here is to automate deployment, having a single action.
In order to deploy `hello` application to production environment we need two actions `bump app version` & `terraform apply`.
```
git checkout <version_to_deploy>
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
