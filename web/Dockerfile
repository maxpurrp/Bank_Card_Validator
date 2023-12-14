FROM python

WORKDIR /web

COPY . .

EXPOSE 3336

RUN python3 -m pip install --no-cache-dir --no-warn-script-location --upgrade pip

RUN pip install django