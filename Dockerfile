FROM docker.test.com/paas/centos:7.5.1804


ADD ./ /srv/riskengine/

ENTRYPOINT  /srv/riskengine/riskengine
