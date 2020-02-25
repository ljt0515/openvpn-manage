client
dev tun
proto {{ .Proto }}
remote {{ .ServerAddress }} {{ .Port }}
resolv-retry infinite
nobind
persist-key
persist-tun
remote-cert-tls server
auth {{ .Auth }}
cipher {{ .Cipher }}
ignore-unknown-option block-outside-dns
block-outside-dns
verb 3
<ca>
-----BEGIN CERTIFICATE-----
{{ .Ca }}
</ca>
<cert>
-----BEGIN CERTIFICATE-----
{{ .Cert }}
</cert>
<key>
-----BEGIN PRIVATE KEY-----
{{ .Key }}
</key>
<tls-crypt>
-----BEGIN OpenVPN Static key V1-----
{{ .TlsCrypt }}
</tls-crypt>