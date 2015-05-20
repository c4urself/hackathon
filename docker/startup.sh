#!/bin/sh

service nginx start
/usr/sbin/sshd
/usr/bin/supervisord -n
