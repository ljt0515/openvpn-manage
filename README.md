# OpenVpn-管理

## 摘要
OpenVPN服务器Web管理界面。

目标：创建易于部署且易于使用的解决方案，使在小型OpenVPN环境中的工作变得轻而易举。

如果您已安装docker和docker-compose，则可以直接跳至[installation](#Prod)。

![Status page](docs/images/preview_status.png?raw=true1)

请注意此项目处于alpha阶段。它仍然需要一些工作，以使它的安全和功能完成。



## 功能

*状态页面，显示服务器统计信息和连接的客户端列表
*轻松创建客户证书
*能够将客户端证书作为ovpn配置文件下载，并带有内部客户端配置
*日志预览
*通过Web界面修改OpenVPN配置文件

## 截图

[截图](docs/screenshots.md)

## Usage

启动后，在端口8080上可以看到web服务。若要登录，请使用下列默认凭据：
* 用户名: admin
* 密  码: 123456 
请立即将密码更改为自己的密码！

### Prod

要求:
* 在防火墙打开的端口上：1194/udp和8080/tcp

执行命令

    curl -O https://raw.githubusercontent.com/adamwalach/openvpn-manage/master/docs/docker-compose.yml
    docker-compose up -d

它会启动两个码头集装箱。一个用于OpenVPN服务器，另一个用于OpenVPNAdmin Web应用程序。通过一个停靠器卷，它创建了以下目录结构：

    .
    ├── docker-compose.yml
    └── openvpn-data
        ├── conf
        │   ├── ca.crt
        │   ├── ca.key
        │   ├── client-common.txt
        │   ├── crl.pem
        │   ├── dh.pem
        │   ├── easy-rsa
        │   ├── ipp.txt
        │   ├── openvpn-status.log
        │   ├── ovpn
        │   ├── server.conf
        │   ├── server.crt
        │   ├── server.key
        │   ├── tc.key
        │   ├── openvpn.log
        │   ├── 
        │   └── 
        └── db
            └── data.db



### Dev

Requirements:
* golang environments
* [beego](https://beego.me/docs/install/)

Execute commands:

    go get openvpn-manage
    cd $GOPATH/src/openvpn-manage
    bee run -gendoc=true

## Todo

* add unit tests
* add option to modify certificate properties
* generate random admin password at initialization phase
* add versioning
* add automatic ssl/tls (check how [ponzu](https://github.com/ponzu-cms/ponzu) did it)


## License

This project uses [MIT license](LICENSE)

## Remarks

### Vendoring
https://github.com/kardianos/govendor is used for vendoring.

To update dependencies from GOPATH:

`govendor update +v`

### Template
AdminLTE - dashboard & control panel theme. Built on top of Bootstrap 3.

Preview: https://almsaeedstudio.com/themes/AdminLTE/index2.html

