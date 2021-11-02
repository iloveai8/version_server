FROM harbor.nuclearport.com/devops/centos:7.9
MAINTAINER Jackpotland
WORKDIR /home/service/game_slots_vsn
COPY game_slots_vsn .
COPY conf/ conf/
ENV ENV dev
RUN pwd \
    && env \
    && ls -al
ENTRYPOINT ./game_slots_vsn  --config=./conf/${ENV}.yaml
