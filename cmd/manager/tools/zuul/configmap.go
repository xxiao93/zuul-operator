package zuul

var mainyamlTemplate = `
- tenant:
  name: es-tenant
  source:
	gerrit:
	  config-projects:
		- easystack/project-config
	  untrusted-projects:
		- easystack/zuul-jobs:
			shadow: easystack/project-config
		- easystack/easystack-zuul-jobs
		- easystack/requirements
		- easystack/neutron
		- easystack/espresso
		- easystack/ark
		- easystack/cube
		- easystack/ark-infra
		- easystack/ems-appcenter-dashboard
		- easystack/rally
		- easystack/tempest
		- easystack/charts
		- easystack/nova
		- easystack/cinder
		- easystack/heat
		- easystack/billing
		- easystack/ceilometer
		- easystack/chakra
		- easystack/gnocchi
		- easystack/ironic
		- easystack/manila
		- easystack/murano
		- easystack/sahara
		- easystack/tickets
		- easystack/trove
		- easystack/estack-hagent
		- easystack/keystone
		- easystack/murano-tempest-plugin
		- easystack/dr
		- easystack/python-drclient
		- easystack/python-muranoclient
		- easystack/python-billingclient
		- easystack/python-chakraclient
		- easystack/python-ticketsclient
		- easystack/python-ceilometerclient
		- easystack/python-cinderclient
		- easystack/python-espressoclient
		- easystack/python-harborclient
		- easystack/python-heatclient
		- easystack/python-keystoneclient
		- easystack/python-magnumclient
		- easystack/python-manilaclient
		- easystack/python-neutronclient
		- easystack/python-novaclient
		- easystack/python-openstackclient
		- easystack/python-troveclient
		- easystack/python-rollerclient
		- easystack/ems-servicecatalog-dashboard
		- easystack/coaster
		- easystack/roller-dashboard
		- easystack/keystoneauth
		- easystack/releases
		- easystack/alertsaver
		- easystack/prometheus-openstack-exporter
		- easystack/eks-dashboard
		- easystack/escmp-dashboard
		- easystack/ecs-dashboard
		- easystack/heat-dashboard
		- easystack/ems-dashboard
		- easystack/dr-dashboard
		- easystack/diamond
		- easystack/peak
		- easystack/ESStorage
		- easystack/django-openstack-auth
`

var zuulconfigTemplate = `
[gearman]
server=127.0.0.1
port=4730
[zookeeper]
hosts=zookeeper-zookeeper:2181
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
`
