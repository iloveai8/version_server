FROM harbor.nuclearport.com/devops/centos:7.9
MAINTAINER Jackpotland
WORKDIR /home/service/gsv
COPY gsv .
COPY conf/ conf/
ENTRYPOINT ./gsv  --config=./conf/${env}.yaml

#docker build --no-cache -t harbor.nuclearport.com/jackpotland/gsv:1.0.0 -f Dockerfile .
#docker run -p 9092:9091 -e env=pro --name gsv1-1.0.0 -d harbor.nuclearport.com/jackpotland/gsv:1.0.0
#docker exec -it gsv1-1.0.0 /bin/bash