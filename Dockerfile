FROM ubuntu:21.10

RUN apt-get update
RUN apt-get install -y openssh-server
RUN apt-get install -y rsync
RUN apt-get install -y software-properties-common
RUN apt-get install -y sshpass

RUN apt-get install -y python3
RUN apt-get install -y python3-pip

RUN pip3 install watchdog
RUN pip3 install flask

COPY . /app
WORKDIR /app
CMD [ "./start" ]