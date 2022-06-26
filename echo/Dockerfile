FROM centos:7.3.1611 

COPY echo /app/echo

RUN yum install -y telnet bind-utils

EXPOSE 8080

WORKDIR /app

CMD ./echo

