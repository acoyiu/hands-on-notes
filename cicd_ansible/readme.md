```sh
docker run -dit --name containerAnsible -v \$PWD:/ansible --net=host ubuntu
```

```sh
docker exec -it containerAnsible bash
apt update
apt install iputils-ping openssh-client software-properties-common
add-apt-repository --yes --update ppa:ansible/ansible
apt-get install -y sshpass
apt install ansible
```

```sh
cd /ansible
ansible-playbook ./1_playbook.yaml -i inventory.ini
```
