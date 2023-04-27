FROM ubuntu

RUN apt update
RUN apt install openssh-server -y
RUN echo 'root:password' | chpasswd
RUN echo "Port 22" >> /etc/ssh/sshd_config
RUN echo "PasswordAuthentication yes" >> /etc/ssh/sshd_config
RUN echo "PermitRootLogin yes" >> /etc/ssh/sshd_config
RUN apt install sshpass

EXPOSE 22

ENTRYPOINT ["/bin/sh", "-c" , "/etc/init.d/ssh restart && tail -f /dev/null"]