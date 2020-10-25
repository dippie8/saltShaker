FROM centos/systemd
RUN yum install epel-release gcc gcc-c++ -y
RUN yum install python3-pip python3-devel -y
RUN yum install -y supervisor
RUN pip3 install pyzmq PyYAML pycrypto msgpack-python jinja2 psutil futures tornado
RUN pip3 install salt==3000.1
ADD salt-minion.service /etc/systemd/system/
ADD salt-master.service /etc/systemd/system/
RUN /usr/sbin/init &
RUN systemctl enable salt-master
RUN systemctl enable salt-minion
RUN mkdir /etc/salt
RUN echo 'auto_accept: True' > /etc/salt/master
CMD ["/usr/sbin/init"]


# docker rm saltshaker -f
# docker build -t saltshaker .
# docker run --privileged --name saltshaker --hostname salt -d saltshaker:latest

# docker rm saltshaker -f && docker build -t saltshaker . && docker run --privileged --name saltshaker --hostname salt -d saltshaker:latest && docker exec -it saltshaker /bin/bash