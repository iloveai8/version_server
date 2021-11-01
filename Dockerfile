FROM harbor.nuclearport.com/devops/centos:7.9
MAINTAINER Jackpotland
WORKDIR /home/service
COPY game_slots_vsn /home/service/
COPY conf/*.yaml /home/service/conf/
ENTRYPOINT /home/service/game_slots_vsn "--config=/home/service/conf/dev.yaml"

