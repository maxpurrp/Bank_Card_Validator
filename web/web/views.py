from django.shortcuts import render, redirect
from django.http import HttpRequest


def web(request):
    return render(request, "main/index.html")


def recall(request: HttpRequest):
    data = request.get_full_path_info().split("?")[1]
    print(data)
    return redirect(f"http://host.docker.internal:3333/check_number?{data}")
