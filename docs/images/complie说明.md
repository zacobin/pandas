一、service compile
1、Build all services.  bin file will be created at directory "bin"
   make    
2、Build specific service
   eg.  make lbs
3、 Build all  services of adaptor
    make adaptor
4、Build all  service of addone
    make addone

二、docker compile and deploy

step1、Build base images
   make pandas-base
step2、Build all pandas services' images
   make dockers_dev
step3、Build all adaptor services' images
   make dockers_adaptor  
step4、Build all addone services' images
   make dockers_addone
Step5、 Deploy by docker-compose
    docker-compose -f docker/docker-compose.yaml up

others：
Build specifice service image
eg. make docker_dev_lbs