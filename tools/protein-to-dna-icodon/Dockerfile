FROM supraja968/my_r_base_image@sha256:13627cc37a9554b86d7562d4853aa31cca5ce01b92ad7d2f2e4f1f1a44fdf24b

RUN apt-get update && apt-get install -y --no-install-recommends \
    python3 \
    python3-pip \
    python3-venv \
    python3-dev \ 
    libblas-dev

WORKDIR /app
RUN python3 -m venv /opt/venv

ENV PATH="/opt/venv/bin:$PATH"

COPY requirements.txt /app/

RUN pip install --no-cache-dir -r requirements.txt

COPY main.py /app
COPY iCodonScript.R /app

ENTRYPOINT ["python", "-u", "main.py"]
