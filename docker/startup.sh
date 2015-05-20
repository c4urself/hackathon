#!/bin/sh

service nginx start
/usr/sbin/sshd
/usr/bin/redis-server
/usr/bin/supervisord -n
