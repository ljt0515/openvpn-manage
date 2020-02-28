FROM debian:jessie
WORKDIR /opt
EXPOSE 8080


ADD assets/start.sh /opt/start.sh
ADD assets/vars.template /opt/scripts/

ADD openvpn-manage.tar.gz /opt/openvpn-gui/
RUN rm -f /opt/openvpn-gui/data.db
ADD assets/app.conf /opt/openvpn-gui/conf/app.conf

CMD /opt/start.sh
