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


### infrastructure
set the following variables in both `var.tf` and `terraform.tfvars`.  
location of variables does simplify any further maintenance of infrastructure.  
but we also have to think about sensitive data.  
the first file is being committed into source and consists of variable parameters.  
the latter file is sensitive, needs to be created manually based on the template below.  
```
cat > terraform.tfvars <<EOF
access_key = "your aws access key"
secret_key = "your aws secret key"
region     = "aws region of operation"
EOF
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
```

### integration test
```
terraform state show aws_lb.alb | grep dns_name
curl <dns_name>:5000/ | jq
```
