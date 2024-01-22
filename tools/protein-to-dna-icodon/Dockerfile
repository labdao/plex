FROM quay.io/labdao/r-base-with-icodon@sha256:dd0b68b6bd07b8830559098b1116b190f04ba46ecfc20866364d1eaea862ab58

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
