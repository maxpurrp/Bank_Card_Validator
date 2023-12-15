from django.shortcuts import render, redirect
from django.http import HttpRequest
import json


def web(request: HttpRequest):
    with open('/opt/q.txt', 'a') as f:
        f.write(str(request.body))
        body = request.body
        res = body.decode()
        f.write(res)
        return render(request, "main/index.html", {'errors': json.loads(res[:-1])})


def recall(request: HttpRequest):
    print('im herer in redirect')
    data = request.get_full_path_info().split("?")[1]
    return redirect(f"http://host.docker.internal:3333/check_number?{data}")
