# for Development
```sh
# test
helm install ingress-hl . --debug --dry-run

# real run
helm install ingress-hl .
```

# for ansible deployment
```sh
ansible-playbook ./ansible/playbook.yaml -i ./ansible/inventory.ini
```