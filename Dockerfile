FROM ubuntu:22.04

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

RUN mkdir /root/.ssh
RUN mv /app/ssh/id_ed25519 /root/.ssh
RUN chown -R root:root /root/.ssh/id_ed25519
RUN chmod 700 /root/.ssh
RUN chmod 600 /root/.ssh/id_ed25519

CMD [ "./start" ]