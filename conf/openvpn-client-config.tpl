client
dev tun
proto {{ .Proto }}
remote {{ .ServerAddress }} {{ .Port }}
resolv-retry infinite
nobind
persist-key
persist-tun
remote-cert-tls server
auth SHA512
cipher AES-256-CBC
ignore-unknown-option block-outside-dns
block-outside-dns
verb 3
<ca>
-----BEGIN CERTIFICATE-----
{{ .Ca }}
-----END CERTIFICATE-----
</ca>
<cert>
-----BEGIN CERTIFICATE-----
{{ .Cert }}
-----END CERTIFICATE-----
</cert>
<key>
-----BEGIN ENCRYPTED PRIVATE KEY-----
{{ .Key }}
-----END ENCRYPTED PRIVATE KEY-----
</key>
<tls-crypt>
-----BEGIN OpenVPN Static key V1-----
{{ .TlsCrypt }}
-----END OpenVPN Static key V1-----
</tls-crypt>