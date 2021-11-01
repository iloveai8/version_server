FROM centos:7.9.2009
MAINTAINER Jackpotland
WORKDIR /home/service/game_slots_vsn
#RUN
#RUN ["executable", "param1", "param2"]
#RUN指令创建的中间镜像会被缓存，并会在下次构建中使用。如果不想使用这些缓存镜像，可以在构建时指定--no-cache参数，如：docker build --no-cache
#ADD conf /conf/${ENV_NAME}.yaml
ADD conf .
ADD game_slots_vsn .
#CMD 构建容器后调用，也就是在容器启动时才进行调用。
ENTRYPOINT  GSVsnServer  "--config=/configs/${ENV_NAME}.yml"

