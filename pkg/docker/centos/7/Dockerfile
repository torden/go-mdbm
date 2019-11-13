FROM centos:7
MAINTAINER Torden <tordengit@gmail.com>

RUN yum install -y yum-plugin-fastestmirror.noarch 
RUN yum install -y make gcc gcc-devel g++ g++-devel vim git openssl-devel gcc-c++ readline-devel cppunit-devel curl wget perl-ExtUtils-MakeMaker doxygen epel-release rpm-build rpmdevtools
RUN wget https://bootstrap.pypa.io/2.6/get-pip.py -O /tmp/get-pip.py
RUN python /tmp/get-pip.py
RUN pip install Sphinx==1.3.6
RUN yum update -y

VOLUME ["/data"]
