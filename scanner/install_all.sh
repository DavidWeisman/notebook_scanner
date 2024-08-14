#!/bin/bash

# Update pip packages
pip install -U spacy

# Download spacy models
python -m spacy download en_core_web_trf
python -m spacy download en_core_web_lg

# Install other Python packages
pip install ipynb-py-convert
pip install lxml_html_clean
pip install tensorflow==2.15.0 flask==3.0.2 waitress==3.0.0 jupyterlab==4.1.1 safety==3.0.1 detect-secrets==1.4.0 presidio-analyzer==2.2.353 whispers==2.2.0 GitPython==3.1.42 boto3==1.34.43 tqdm==4.66.3 flask_restx==1.3.0 flask_cors==4.0.0

# Install nbdefense
pip install nbdefense

# Install Trivy using curl and sh
curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /Users/david/Desktop/anaconda3/envs/testenv/lib/python3.9/site-packages/nbdefense/ v0.32.1

