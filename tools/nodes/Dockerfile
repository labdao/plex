FROM nytimes/blender@sha256:194c14fc87a3fec7981f9837edfb3ff19cfc46c7da9392146446d15ee27657ea

RUN wget https://github.com/BradyAJohnston/MolecularNodes/releases/download/v2.6.0/MolecularNodes_2.6.0.zip

COPY . .
RUN blender --background --python ./app.py