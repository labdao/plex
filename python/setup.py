from setuptools import setup, find_packages
from setuptools.command.install import install
from setuptools.command.develop import develop
import os
import subprocess
import shutil
import tempfile
import sys


class PostSetupCommand:
    """Base class for both installation and development mode."""
    def run_setup(self):
        print("Running post-installation...")

        # Retrieve platform from environment variable
        plat_name = os.environ['PLAT_NAME']

        # map plat_name to go_bin_url
        urls = {
            "darwin_x86_64": "https://github.com/labdao/plex/releases/download/v0.8.1/plex_0.8.1_darwin_amd64.tar.gz",
            "darwin_arm64": "https://github.com/labdao/plex/releases/download/v0.8.1/plex_0.8.1_darwin_arm64.tar.gz",
            "linux_x86_64": "https://github.com/labdao/plex/releases/download/v0.8.1/plex_0.8.1_linux_amd64.tar.gz",
            "win_amd64": "https://github.com/labdao/plex/releases/download/v0.8.1/plex_0.8.1_windows_amd64.tar.gz",
        }

        go_bin_url = urls.get(plat_name)

        if go_bin_url:
            try:
                with tempfile.TemporaryDirectory() as temp_dir:
                    self.download_and_extract(go_bin_url, temp_dir)

                    # If self.install_scripts is None, use a default scripts directory
                    if self.install_scripts is None:
                        self.install_scripts = os.path.join(sys.prefix, 'bin')

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


class PostInstallCommand(PostSetupCommand, install):
    """Post-installation for installation mode."""
    def run(self):
        install.run(self)
        self.run_setup()


class PostDevelopCommand(PostSetupCommand, develop):
    """Post-installation for development mode."""
    def run(self):
        develop.run(self)
        self.run_setup()

setup(
    name="PlexLabExchange",
    version="0.8.11",
    packages=find_packages(where='src'),  # tell setuptools to look in the 'src' directory for packages
    package_dir={'': 'src'},  # tell setuptools that the packages are under the 'src' directory
    cmdclass={
        'install': PostInstallCommand,
        'develop': PostDevelopCommand,
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
