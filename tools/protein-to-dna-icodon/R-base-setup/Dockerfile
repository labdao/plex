FROM r-base:latest

RUN apt-get update && apt-get install -y \
    build-essential \
    libcurl4-openssl-dev \
    libssl-dev \
    libxml2-dev \
    libgit2-dev \
    libfontconfig1-dev \
    libharfbuzz-dev \
    libfribidi-dev \
    libpng-dev \
    libtiff5-dev \
    libjpeg-dev \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

RUN apt-get update && apt-get install -f

RUN R -e "options(repos = list(CRAN = 'http://cran.rstudio.com')); \
    install.packages(c('systemfonts', 'textshaping', 'ragg', 'pkgdown', 'devtools'), dependencies = TRUE, verbose = TRUE)"

RUN R -e "devtools::install_github('santiago1234/iCodon')"