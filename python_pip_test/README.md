## Direction to publish on PyPi

# Go to the python_pip_test directory
* `cd python_pip_test`

# Build and publish for every OS
* `python setup.py bdist_wheel --plat-name macosx_10_9_x86_64`
* `twine upload dist/*`

* `python setup.py bdist_wheel --plat-name macosx_11_0_arm64`
* `twine upload dist/*`

* `python setup.py bdist_wheel --plat-name manylinux2014_x86_64`
* `twine upload dist/*`

* `python setup.py bdist_wheel --plat-name win_amd64`
* `twine upload dist/*`
