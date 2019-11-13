FROM ubuntu:16.04
MAINTAINER Torden <tordengit@gmail.com>

RUN sed -e 's/archive/kr.archive/gi' /etc/apt/sources.list > /etc/apt/sources.list.new
RUN mv /etc/apt/sources.list.new /etc/apt/sources.list
RUN apt-get update -y 
RUN apt-get install vim apt-file git curl wget build-essential libssl1.0.0 libssl-dev libreadline6 libreadline6-dbg libreadline6-dev libcppunit-dev libxml2 libxml2-dbg libxml2-dev libxml2-doc libxml2-utils -y
RUN apt-file update

VOLUME ["/data"]
