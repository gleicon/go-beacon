# Elasticsearch, Kibana, Logstash, Ansible and Vagrant

## Depends
    - Vagrant (local)
    - Ansible
    - Python
    - VirtualBox (local)

## Software
    - Ubuntu 14.04 64bits
    - Elastic search 1.1.4
    - Kibana
    - Logstash
    - Nginx
    - Python libs
    - Build tools


# Local
    - vagrant up
    - vagrant ssh 
    - http://192.168.33.22/ - external kibana
    - http://192.169.33.22:8080

# Install on VPS/Cloud/etc
    
    - Adicionar o hostname no arquivo de inventario do Ansible (/etc/ansible/hosts ou local)
        [elasticsearch]
        elasticsearch ansible_ssh_host=192.168.1.1

    - Check passwordless ssh conn
        ssh root@192.168.1.1

    - Run 
        ansible-playbook -l elasticsearch elasticsearch.yml

