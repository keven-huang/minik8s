# serverless_server.py
from flask import Flask, request
import importlib
import schedule
import threading
import time
import requests
import os
import sys

app = Flask(__name__)
request_count = 0
last_request_time = time.time()

FAILED_TIME = 30
Request_20_second = 4

Function_name = ""


# 打开日志文件
script_dir = os.path.dirname(os.path.abspath(__file__))
log_file_path = os.path.join(script_dir, "serverless_server.log")
log_file = open(log_file_path, "a")

# 重定向标准输出到日志文件
sys.stdout = log_file


def reset_timer():
    global last_request_time
    last_request_time = time.time()


def check_timer():
    global last_request_time
    print(time.time(), last_request_time)
    if time.time() - last_request_time >= FAILED_TIME:
        print("No requests received for 1 minute. Exiting...")
        # 进行异常退出的操作，例如抛出异常或者调用系统退出函数
        os._exit(666)  # 退出整个程序

def send_notification(url):
    requests.post(url, json={"message": "Exceeded request limit", "function_name": Function_name})


count_lock = threading.Lock()


def reset_count():
    global request_count
    with count_lock:
        request_count = 0


def increment_count():
    global request_count
    with count_lock:
        request_count += 1

@app.route('/function/<string:module_name>/<string:function_name>', methods=['POST'])
def execute_function(module_name: str, function_name: str):
    global Function_name
    Function_name = function_name
    global request_count
    reset_timer()
    increment_count()

    if request_count >= Request_20_second:
        send_notification("https://192.168.1.7:8080/scale")  # 替换为你要发送通知的URL

    module = importlib.import_module(module_name)
    event = {"method": "http"}

    if request.headers.get('Content-Type') == 'application/json':
        context = request.get_json()
    else:
        context = request.form.to_dict()

    result = getattr(module, function_name)(event, context)

    return result, 200


# 创建一个新线程并启动该线程
def run_thread():
    while True:
        schedule.run_pending()

if __name__ == '__main__':
    job1 = schedule.every(2).seconds.do(check_timer)

    job2 = schedule.every(20).seconds.do(reset_count)

    t = threading.Thread(target=run_thread)
    t.start()

    app.run(host='0.0.0.0', port=8888, processes=True)
