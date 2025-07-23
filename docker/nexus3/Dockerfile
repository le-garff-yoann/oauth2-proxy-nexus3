FROM sonatype/nexus3:3.82.0-alpine

USER root

RUN echo "nexus.scripts.allowCreation=true" >> /opt/sonatype/nexus/etc/nexus-default.properties

USER nexus
