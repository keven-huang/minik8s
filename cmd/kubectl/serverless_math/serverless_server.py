# serverless_server.py
from flask import Flask, request
import importlib
import threading
import time
import requests
import os

app = Flask(__name__)
request_count = 0
last_request_time = time.time()

FAILED_TIME = 30
Request_20_second = 4

Function_name = ""


def reset_timer():
    global last_request_time
    last_request_time = time.time()


def check_timer():
    global last_request_time
    if time.time() - last_request_time >= FAILED_TIME:
        print("No requests received for 1 minute. Exiting...")
        # 进行异常退出的操作，例如抛出异常或者调用系统退出函数
        os._exit(0)  # 退出整个程序

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


if __name__ == '__main__':
    timer_thread = threading.Timer(FAILED_TIME, check_timer)
    timer_thread.start()

    reset_count_thread = threading.Timer(20, reset_count)
    reset_count_thread.start()

    app.run(host='0.0.0.0', port=8888, processes=True)
