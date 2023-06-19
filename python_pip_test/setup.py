from setuptools import setup, find_packages
from setuptools.command.install import install
import os
import tarfile
import platform
import requests


class PostInstallCommand(install):
    """Post-installation for installation mode."""
    def run(self):
        # Call parent 
        print("Running post installation script...")
        install.run(self)
        print("Running post installation script after parent instal")
        # Determine platform 
        go_bin_url = ""
        system_platform = platform.system()
        print(system_platform)
        machine = platform.machine()
        print(machine)
        if system_platform == "Windows":
            if machine == "AMD64":
                go_bin_url = "https://github.com/labdao/plex/releases/download/v0.8.0/plex_0.8.0_windows_amd64.tar.gz"
            elif machine == "i386":
                go_bin_url = "https://github.com/labdao/plex/releases/download/v0.8.0/plex_0.8.0_windows_386.tar.gz"
            elif machine == "ARM64":
                go_bin_url = "https://github.com/labdao/plex/releases/download/v0.8.0/plex_0.8.0_windows_arm64.tar.gz"
        elif system_platform == "Linux":
            if machine == "x86_64":
                go_bin_url = "https://github.com/labdao/plex/releases/download/v0.8.0/plex_0.8.0_linux_amd64.tar.gz"
            elif machine == "i386":
                go_bin_url = "https://github.com/labdao/plex/releases/download/v0.8.0/plex_0.8.0_linux_386.tar.gz"
            elif machine == "aarch64":
                go_bin_url = "https://github.com/labdao/plex/releases/download/v0.8.0/plex_0.8.0_linux_arm64.tar.gz"
        elif system_platform == "Darwin":
            if machine == "x86_64":
                go_bin_url = "https://github.com/labdao/plex/releases/download/v0.8.0/plex_0.8.0_darwin_amd64.tar.gz"
            elif machine == "arm64":
                go_bin_url = "https://github.com/labdao/plex/releases/download/v0.8.0/plex_0.8.0_darwin_arm64.tar.gz"
        # Download Go binary to scripts path
        if go_bin_url:
            try:
                r = requests.get(go_bin_url, allow_redirects=True)
                r.raise_for_status()  # Will raise an exception if the status code is not 200
            except requests.exceptions.RequestException as e:
                raise RuntimeError(f"Failed to download the Go binary: {e}")
            # download the tar file
            tar_file_path = os.path.join(self.install_lib, "plex.tar.gz")
            try:
                with open(tar_file_path, 'wb') as f:
                    f.write(r.content)
            except Exception as e:
                raise RuntimeError(f"Failed to write the Go binary to disk: {e}")
            # extract the binary from the tar file
            try:
                with tarfile.open(tar_file_path, 'r:gz') as tar:
                    tar.extractall(path=self.install_scripts)
            except tarfile.TarError as e:
                raise RuntimeError(f"Failed to extract the Go binary: {e}")
            # set the binary as executable
            dst = os.path.join(self.install_scripts, "plex")
            try:
                os.chmod(dst, 0o755)  # make sure the binary is executable
            except Exception as e:
                raise RuntimeError(f"Failed to set the Go binary as executable: {e}")
        else:
            raise RuntimeError(f"The current platform/machine of {system_platform}/{machine} is not supported.")


setup(
    name="pyplex",
    version="0.8.0",
    packages=find_packages(),
    cmdclass={
        'install': PostInstallCommand,
    },
    install_requires=[
        'requests',
    ],
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
