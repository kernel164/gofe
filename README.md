## gofe - Go File Explorer
A golang backend for angular-filemanager - https://github.com/joni2back/angular-filemanager

### Screenshots
![](https://raw.githubusercontent.com/kernel164/gofe/master/screenshot1.png)
![](https://raw.githubusercontent.com/kernel164/gofe/master/screenshot2.png)

### Features
- Login support
- SSH backend support

### Sample Config
```ini
BACKEND = ssh
SERVER = http

[server.http]
BIND = localhost:4000
STATICS = angular-filemanager/bower_components,angular-filemanager/dist,angular-filemanager/src

[backend.ssh]
HOST = localhost:22
```