# Overview
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fcloustone%2Fpandas.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fcloustone%2Fpandas?ref=badge_shield) ![Go](https://img.shields.io/badge/golang-1.13.8-green.svg) ![Go](https://github.com/cloustone/pandas/workflows/Go/badge.svg) 

Pandas is device managment and modeling SaaS which can be deployed in any environments to provide the following features.

* Integrate with 3rd-party broker, for example, mainflux to provide device connection. 
* Device management
* Device model management, user can creating device model using web-frontend.
* Project management,
* Full featured IoT rule engine
* Location Based Service
* Display project's realtime device status based on device model
* SCADA
* Security subsystem
* Deployment with docker-compose and Kubernetes

Pandas is designed and implemented based on micro service architecture,which include the following components.

* Dashboard, the web console to manage all objects(User/Project/Workshop/Device/RuleChain).
* ApiMachinery, the API gateway.
* Dmms: Device model management service.
* Pms: Project management service.
* Rulechain: Rule chain service
* Headmast: A simple job scheduler
* Shrio: Security Manager

![](docs/images/pandas-arch.png)

### Depencency
* go-swagger
* go-bindata 
* protoc-gen-go

### Building

````
git clone --recursive https://github.com/cloustone/pandas .git

Verify the dashboard can be build rightly
> cd $GOPATH/src/github.com/cloustone/pandas/dashboard
> npm install
> npm run dev

Build dashboard into go-bindata
>   cd ..
>  ./scripts/dashboard_gen.sh

build and run 
> cd $GOPATH/src/github.com/cloustone/pandas  
> make 
> make docker
> docker-compose up -d
````

## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fcloustone%2Fpandas.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fcloustone%2Fpandas?ref=badge_large)
