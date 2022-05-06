FROM harbor.nuclearport.com/devops/centos:7.9
MAINTAINER Jackpotland
WORKDIR /home/service/gsv
COPY gsv .
COPY conf/ conf/
COPY static/ static/
ENTRYPOINT ./gsv  start --config=./conf/${env}.yaml

#docker build --rm --no-cache -t harbor.nuclearport.com/jackpotland/gsv:dev -f Dockerfile .
#docker run --rm -p 7001:9091 -e env=dev --name gsv-dev -d harbor.nuclearport.com/jackpotland/gsv:dev
#docker exec -it gsv-dev /bin/bash

#docker build --rm --no-cache -t harbor.nuclearport.com/jackpotland/gsv:pre -f Dockerfile .
#docker run --rm -p 7002:9091 -e env=pre  --name gsv-pre -d harbor.nuclearport.com/jackpotland/gsv:pre
#docker exec -it gsv-pre /bin/bash

#docker build --rm --no-cache -t harbor.nuclearport.com/jackpotland/gsv:pro.1.0.0 -f Dockerfile .
#docker run --rm -p 7003:9091 -e env=pro --name gsv-pro-1.0.0 -d harbor.nuclearport.com/jackpotland/gsv:pro.1.0.0
#docker exec -it gsv-pro-1.0.0 /bin/bash