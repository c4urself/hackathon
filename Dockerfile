FROM ubuntu:14.04
MAINTAINER Ryuk The Shinigami "dev@tubularlabs.com"

# Install base
RUN apt-get update --fix-missing
RUN apt-get install -y git libxslt-dev curl libcurl4-openssl-dev software-properties-common python-software-properties nginx libev-dev golang

# Enable SSH
RUN apt-get install -y openssh-server
RUN echo 'root:hackathon' | chpasswd
RUN mkdir /var/run/sshd

# Install Nginx
RUN add-apt-repository ppa:nginx/stable -y
RUN apt-get update
RUN apt-get install -y nginx

ADD . /opt/hackathon
WORKDIR /opt/hackathon
ENV GOPATH /opt/hackathon
RUN go get ./...

# Setup nginx and supervisor
RUN rm /etc/nginx/sites-enabled/default

RUN cp /opt/hackathon/docker/supervisord.conf /etc/supervisor/conf.d/hackathon.conf
RUN cp /opt/hackathon/nginx.conf /etc/nginx/sites-available/hackathon.conf
RUN ln -s ../sites-available/hackathon.conf /etc/nginx/sites-enabled/hackathon.conf
RUN cp /opt/hackathon/docker/startup.sh /opt/startup.sh

EXPOSE 22 80 8080

RUN chmod 755 /opt/startup.sh
CMD ["/bin/sh", "/opt/startup.sh"]
