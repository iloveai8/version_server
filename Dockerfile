FROM alpine
ADD aman /aman
ADD configs /configs
ENV ENV_NAME test
ENTRYPOINT  /aman "--config=/configs/${ENV_NAME}.yml"

