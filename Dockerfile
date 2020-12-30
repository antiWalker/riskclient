FROM docker.mfwdev.com/paas/centos:7.5.1804

# 更新 dockerfile 地址 https://aos.mfwdev.com/microservice/updateapp
# appid 3429

ADD ./ /srv/riskengine/

ENTRYPOINT  /srv/riskengine/riskengine
