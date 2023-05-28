FROM python:3.10-alpine
MAINTAINER LuHaoqi <LuHaoqi123@sjtu.edu.cn>

# 更新操作系统和安装依赖
RUN apk update && apk upgrade && \
    apk add --no-cache bash

# install python requirements
RUN pip install flask importlib requests request

# 如果用户需要安装额外的库，那么还需要提供requirements.txt
RUN pip install --no-cache-dir -r requirements.txt

# 把两个文件添加到image中，这样就不需要挂载了
COPY ./serverless_server.py serverless_server.py
COPY ./my_module.py my_module.py

# define command for this image
CMD ["python3","/serverless_server.py"]