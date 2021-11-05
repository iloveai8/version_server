FROM harbor.nuclearport.com/devops/centos:7.9
MAINTAINER Jackpotland
WORKDIR /home/service/gsv
COPY gsv .
COPY conf/ conf/
ENV env dev
ENV vsn 1.0.0
RUN pwd \
    && env \
    && ls -al
ENTRYPOINT ./gsv  --config=./conf/${env}.yaml
