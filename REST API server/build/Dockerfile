## Skeleton Dockerfile for your service
## Change folder name: <yaaf> to your product folder name
## Change executable name: <yaaf-api-server> to your server executable name
## Add additional apt packages in required
##
FROM ubuntu:22.04

RUN apt-get update
RUN apt-get -yq install apt-utils
RUN apt-get -yq install ca-certificates && update-ca-certificates


RUN mkdir yaaf && cd yaaf

ADD yaaf-api-server /yaaf/.
ADD build-tag /yaaf/.

RUN chmod +x /yaaf/yaaf-api-server

WORKDIR "/yaaf"

CMD ["./yaaf-api-server"]
