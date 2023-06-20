Direction to publish on PyPi
From the python directory
On Linux
* `python3 setup.py bdist_wheel --plat-name manylinux1_x86_64`

On Mac M1
* `python setup.py bdist_wheel --plat-name macosx_11_0_arm64`

On Windows
* TODO

On Mac Intel
* TODO

Then run to publish
* `twine upload dist/*`
