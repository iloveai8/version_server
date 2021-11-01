FROM harbor.nuclearport.com/devops/centos:7.9
MAINTAINER Jackpotland
WORKDIR /home/service

ADD game_slots_vsn /game_slots_vsn
ADD conf /conf

ENV ENVNAME dev
ENTRYPOINT ["game_slots_vsn", "--config=/conf/${ENV_NAME}.yml"]

