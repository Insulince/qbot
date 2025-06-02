@echo off
set dt=%DATE:~4%
set dt=%dt:/=-%
set dt=%dt: =%
set tm=%TIME::=-%
set tm=%tm: =%
set timestamp=%dt%_%tm:~0,8%
flyctl ssh sftp get /var/lib/litefs/qbot.db ./backups/%timestamp%.sqlite
