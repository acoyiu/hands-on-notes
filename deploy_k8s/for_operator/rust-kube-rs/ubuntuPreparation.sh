# # create Docker container
# docker run -dit --name rust -v $PWD:/app -v /home/aco/.kube/config:/root/.kube/config ubuntu sh -c 'sleep infinity'
# docker exec -it rust bash

# install package
apt update
apt install curl libssl-dev pkg-config build-essential -y

# install rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# Add this to the top of file, /etc/ssl/openssl.cnf
echo -e "openssl_conf = default_conf\n$(cat /etc/ssl/openssl.cnf)" > /etc/ssl/openssl.cnf
source "$HOME/.cargo/env"

# Add this to the bottom of file
cat >> /etc/ssl/openssl.cnf << EOF
[ default_conf ]
ssl_conf = ssl_sect
[ssl_sect]
system_default = system_default_sect
[system_default_sect]
MinProtocol = TLSv1.2
CipherString = DEFAULT:@SECLEVEL=1
EOF