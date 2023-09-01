FROM us-docker.pkg.dev/colab-images/public/runtime
RUN apt update -y \
    && apt install -y wget nano
# RUN python3 -m pip uninstall -y jupyter jupyter_core jupyter-client jupyter-console jupyterlab_pygments qtconsole notebook nbconvert nbformat nbclassic nbclient jupyterlab-widgets jupyter-events jupyter-server jupyter-server-terminals 

RUN pip3 install hydra-core --upgrade
RUN mkdir /app
COPY install.py /app
RUN cd app && python3 install.py
COPY . /app
ENTRYPOINT ["/bin/bash"]