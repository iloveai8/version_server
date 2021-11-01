FROM harbor.nuclearport.com/devops/centos:7.9
MAINTAINER Jackpotland
#WORKDIR /home/service
#RUN
#RUN ["executable", "param1", "param2"]
#RUN指令创建的中间镜像会被缓存，并会在下次构建中使用。如果不想使用这些缓存镜像，可以在构建时指定--no-cache参数，如：docker build --no-cache
#COPY conf /home/service
#COPY game_slots_vsn /home/service

ADD game_slots_vsn /game_slots_vsn
ADD conf /conf

ENV ENVNAME dev
ENTRYPOINT /game_slots_vsn  "--config=/conf/${ENV_NAME}.yml"

