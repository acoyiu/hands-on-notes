FROM ubuntu

RUN apt update
RUN apt-get install -y software-properties-common
RUN apt-get install -y openssh-client
RUN apt-get install -y sshpass 
RUN add-apt-repository --yes --update ppa:ansible/ansible
RUN apt install -y ansible

ENTRYPOINT ["tail", "-f", "/dev/null"]