FROM python:3.10-slim

# Instala las bibliotecas necesarias
RUN apt-get update && apt-get install -y \
    gcc \
    cron \
    libpq-dev \
    && rm -rf /var/lib/apt/lists/*


RUN ln -s libcrypto.so.1.1 /usr/local/lib/libcrypto.so
RUN apt-get update
RUN apt install -y libpq-dev python3-dev 

RUN mkdir /app
WORKDIR /app


COPY requirements.txt .
RUN pip install -r requirements.txt

COPY tasks /app/tasks
RUN ls -l /app/tasks

RUN chmod +x /app/tasks

COPY tasks/backup/cronjob /etc/cron.d/cronjob
RUN chmod 0644 /etc/cron.d/cronjob 
RUN touch /var/log/cron.log
RUN crontab /etc/cron.d/cronjob

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

CMD ["cron", "-f"]
