from setuptools import setup, find_packages

setup(
    name="mypackage",
    version="0.1.0",
    packages=find_packages(),
    install_requires=[
        "pydantic==1.8.2",
        "biopandas==0.2.9",
        "rdkit==2021.03.5",
        "pandas==1.3.3",
    ],
)
