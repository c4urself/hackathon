FROM ubuntu:14.04
MAINTAINER Ryuk The Shinigami "dev@tubularlabs.com"

# Install base
RUN apt-get update --fix-missing
RUN apt-get install -y git libxslt-dev build-essential libcurl4-openssl-dev software-properties-common python-software-properties nginx libev-dev supervisor

# Enable SSH
RUN apt-get install -y openssh-server
RUN echo 'root:tubular' | chpasswd
RUN sed -i '/PermitRootLogin/s/without-password/yes/' /etc/ssh/sshd_config
RUN sed 's@session\s*required\s*pam_loginuid.so@session optional pam_loginuid.so@g' -i /etc/pam.d/sshd
RUN mkdir /var/run/sshd

# Install Nginx
RUN add-apt-repository ppa:nginx/stable -y
RUN apt-get update
RUN apt-get install -y nginx

# Install redis
RUN apt-get -y install redis-server

# Install latest Golang
RUN wget https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.4.2.linux-amd64.tar.gz
ENV PATH /usr/local/go/bin:$PATH
RUN go version

RUN mkdir -p /opt/go/src/github.com/c4urself
ADD . /opt/go/src/github.com/c4urself/hackathon
WORKDIR /opt/go/src/github.com/c4urself/hackathon
ENV GOPATH /opt/go
RUN make vendor_update
RUN make

# Setup nginx and supervisor
RUN rm /etc/nginx/sites-enabled/default
RUN cp /opt/go/src/github.com/c4urself/hackathon/docker/supervisord.conf /etc/supervisor/conf.d/hackathon.conf
RUN cp /opt/go/src/github.com/c4urself/hackathon/docker/nginx.conf /etc/nginx/sites-available/hackathon.conf
RUN ln -s ../sites-available/hackathon.conf /etc/nginx/sites-enabled/hackathon.conf
RUN cp /opt/go/src/github.com/c4urself/hackathon/docker/startup.sh /opt/startup.sh

EXPOSE 22 80 8080

RUN chmod 755 /opt/startup.sh
CMD ["/bin/sh", "/opt/startup.sh"]
