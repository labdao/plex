from setuptools import setup, find_packages
from setuptools.command.install import install
import os
import platform
import subprocess
import shutil
import tempfile


class PostInstallCommand(install):
    """Post-installation for installation mode."""
    def run(self):
        install.run(self)

        print("Running post-installation...")

        # Determine platform 
        os_type = platform.system().lower()
        arch = platform.machine().lower()

        # set default go bin url 
        go_bin_url = None

        if os_type == "darwin" and arch in ["amd64", "x86_64", "i386"]:
            go_bin_url = "https://github.com/labdao/plex/releases/download/v0.8.0/plex_0.8.0_darwin_amd64.tar.gz"
        elif os_type == "darwin" and arch == "arm64":
            go_bin_url = "https://github.com/labdao/plex/releases/download/v0.8.0/plex_0.8.0_darwin_arm64.tar.gz"
        elif os_type == "linux" and arch in ["amd64", "x86_64"]:
            go_bin_url = "https://github.com/labdao/plex/releases/download/v0.8.0/plex_0.8.0_linux_amd64.tar.gz"
        elif os_type == "windows" and arch in ["amd64", "x86_64"]:
            go_bin_url = "https://github.com/labdao/plex/releases/download/v0.8.0/plex_0.8.0_windows_amd64.tar.gz"

        if go_bin_url:
            try:
                with tempfile.TemporaryDirectory() as temp_dir:
                    self.download_and_extract(go_bin_url, temp_dir)

                    # move the binary to the scripts installation directory
                    src = os.path.join(temp_dir, 'plex')
                    dst = os.path.join(self.install_scripts, 'plex')

                    # Create target Directory if don't exist
                    if not os.path.exists(self.install_scripts):
                        os.makedirs(self.install_scripts)

                    shutil.move(src, dst)
                    # set the binary as executable
                    os.chmod(dst, 0o755)

            except Exception as e:
                print(f"Failed to download and extract the Go binary: {str(e)}")
                raise

    def download_and_extract(self, go_bin_url, temp_dir):
        subprocess.run(f"curl -sSL {go_bin_url} | tar xvz -C {temp_dir}", shell=True, check=True)


setup(
    name="PlexLabExchange",
    version="0.8.6",
    packages=find_packages(),
    cmdclass={
        'install': PostInstallCommand,
    },
    author="LabDAO",
    author_email="media@labdao.xyz",
    description="A Python interface to the Plex Go CLI.",
    url="https://github.com/labdao/plex",
    classifiers=[
        "Development Status :: 3 - Alpha",
        "Intended Audience :: Developers",
        "License :: OSI Approved :: MIT License",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
    ],
    keywords="plex golang cli wrapper",
    license="MIT",
    python_requires='>=3.8',
)
