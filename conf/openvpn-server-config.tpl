management {{ .Management }}

local 192.168.10.155
port {{ .Port }}
proto {{ .Proto }}
dev tun
ca {{ .Ca }}
cert {{ .Cert }}
key {{ .Key }}
max-clients {{ .MaxClients }}
dh {{ .Dh }}
auth {{ .Auth }}
tls-crypt tc.key
topology subnet
server {{ .Server }}
ifconfig-pool-persist {{ .IfconfigPoolPersist }}
push "redirect-gateway def1 bypass-dhcp"
push "dhcp-option DNS 114.114.114.114"
push "dhcp-option DNS 8.8.8.8"
keepalive {{ .Keepalive }}
cipher {{ .Cipher }}
user nobody
group nobody
persist-key
persist-tun
status openvpn-status.log
log openvpn.log
verb 3
crl-verify crl.pem
explicit-exit-notify
