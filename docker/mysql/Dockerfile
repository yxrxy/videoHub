FROM mysql:8.0

ENV MYSQL_ROOT_PASSWORD=videoHub
ENV MYSQL_DATABASE=videoHub

# 将初始化SQL复制到docker-entrypoint-initdb.d目录
# MySQL容器首次启动时会自动执行这个目录下的.sql文件
COPY ./init.sql /docker-entrypoint-initdb.d/

ENV TZ=Asia/Shanghai
# MySQL配置
COPY ./mysql.cnf /etc/mysql/conf.d/

EXPOSE 3306 