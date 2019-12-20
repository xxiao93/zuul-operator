package zuul

var zuul_user_id int64 = 42488

var mainyamlTemplate = `
- tenant:
    name: es-tenant
    source:
      gerrit:
        config-projects:
          - xxxxx/project-config
        untrusted-projects:
          - xxxxx/zuul-jobs:
              shadow: xxxxx/project-config
          - xxxxx/zuul-jobs
          - xxxxx/requirements
          - xxxxx/neutron
`

var zuulschedulerconfigTemplate = `
[gearman]
server=127.0.0.1
port=4730
[zookeeper]
hosts=zookeeper.devops.svc.cluster.local:2181
[gearman_server]
start=true
port=4730
listen_address=0.0.0.0
[scheduler]
tenant_config=/var/lib/zuul/tenant-config/main.yaml
pidfile=/var/lib/zuul/run/zuul-scheduler.pid
[connection gerrit]
driver=gerrit
server={{ .GerritServer }}
baseurl=http://{{ .GerritServer }}:8080
user={{ .GerritUser }}
sshkey=/home/zuul/.ssh/id_rsa
keepalive=5
[connection mysql]
driver=sql
dburi=mysql+pymysql://zuul:zuul@mysql/zuul
`

var zuulexecutorconfigTemplate = `
[executor]
user=zuul
finger_port=7900
pidfile=/var/lib/zuul/run/zuul-executor.pid
trusted_ro_paths=/home/zuul/.ssh/id_rsa:/home/zuul/.ssh/id_rsa.pub:/home/zuul/.ssh/known_hosts:/home/zuul/pip.conf:/var/lib/zuul/run/zuul-scheduler.pid
trusted_rw_paths=/var/lib/zuul/tenant-config/
variables=/var/lib/zuul/site-variables.yaml
disk_limit_per_job=4096
[gearman]
server=gearman
port=4730
[connection gerrit]
driver=gerrit
server={{ .GerritServer }}
baseurl=http://{{ .GerritServer }}:8080
user={{ .GerritUser }}
sshkey=/home/zuul/.ssh/id_rsa
keepalive=5
`

var zuulmergerconfigTemplate = `
[merger]
git_dir=/var/lib/zuul/git
git_user_email=zuul@gerrit.com
git_user_name=esadmin
pidfile=/var/lib/zuul/run/zuul-merger.pid
[gearman]
server=gearman
port=4730
[connection gerrit]
driver=gerrit
server={{ .GerritServer }}
baseurl=http://{{ .GerritServer }}:8080
user={{ .GerritUser }}
sshkey=/home/zuul/.ssh/id_rsa
keepalive=5
`

var zuulwebconfigTemplate = `
[web]
listen_address=0.0.0.0
pidfile=/var/lib/zuul/run/zuul-web.pid
port=9001
status_url=http://zuul-web:9001/status
static_path=/zuul/zuul/web/static
[gearman]
server=gearman
port=4730
[zookeeper]
hosts=zookeeper.devops.svc.cluster.local:2181
[connection mysql]
driver=sql
dburi=mysql+pymysql://zuul:zuul@mysql/zuul
`
