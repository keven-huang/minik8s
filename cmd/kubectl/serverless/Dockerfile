FROM ubuntu
MAINTAINER LuHaoqi <LuHaoqi123@sjtu.edu.cn>

# update ubuntu repository
RUN DEBIAN_FRONTEND=noninteractive apt-get -y update

# install ubuntu packages
RUN DEBIAN_FRONTEND=noninteractive apt-get -y install python3 python3-pip

# 如果用户需要安装额外的库，那么还需要提供requirements.txt
RUN pip install -r requirements.txt

# install python requirements
RUN pip install flask

# 把两个文件添加到image中，这样就不需要挂载了
COPY ./serverless_server.py serverless_server.py
COPY ./my_module.py my_module.py

# define command for this image
CMD ["python3","/serverless_server.py"]