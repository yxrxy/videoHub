FROM mysql:8.0

# 设置时区
ENV TZ=Asia/Shanghai

# MySQL配置
COPY ./mysql.cnf /etc/mysql/conf.d/

# 复制初始化脚本
COPY ./init.sql /docker-entrypoint-initdb.d/

EXPOSE 3306 