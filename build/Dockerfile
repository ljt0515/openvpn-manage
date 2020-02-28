FROM centos
WORKDIR /opt
EXPOSE 8080


ADD assets/start.sh /opt/start.sh
ADD assets/vars.template /opt/scripts/

ADD openvpn-manage.tar.gz /opt/openvpn-manage/
RUN rm -f /opt/openvpn-gui/data.db
ADD assets/app.conf /opt/openvpn-manage/conf/app.conf

CMD /opt/start.sh
