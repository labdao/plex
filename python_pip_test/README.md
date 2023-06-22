## Direction to publish on PyPi

# Go to the python_pip_test directory
* `cd python_pip_test`

# Build and publish for every OS
```
export PLAT_NAME=darwin_x86_64
python setup.py bdist_wheel --plat-name macosx_10_9_x86_64
export PLAT_NAME=darwin_arm64
python setup.py bdist_wheel --plat-name macosx_11_0_arm64
export PLAT_NAME=linux_x86_64
python setup.py bdist_wheel --plat-name manylinux2014_x86_64
export PLAT_NAME=win_amd64
python setup.py bdist_wheel --plat-name win_amd64
twine upload dist/*
```
