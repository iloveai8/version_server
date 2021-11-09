FROM harbor.nuclearport.com/devops/centos:7.9
MAINTAINER Jackpotland
WORKDIR /home/service/gsv
COPY gsv .
COPY conf/ conf/

CMD "--config=./conf/${env}.yaml"
ENTRYPOINT ./gsv
