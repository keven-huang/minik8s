FROM ubuntu:latest

RUN apt-get update

RUN DEBIAN_FRONTEND=noninteractive apt-get -y install expect openssh-client

ADD ./upload.sh /upload.sh

RUN chmod +x upload.sh

CMD [ "./upload.sh" ]