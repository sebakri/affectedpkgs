#!/bin/bash
set -e

# Configuration
REPO="sebakri/affectedpkgs"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="affectedpkgs"

# Detect OS
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "${OS}" in
    linux*)  OS='linux';;
    darwin*) OS='darwin';;
    msys*|cygwin*|mingw*) OS='windows';;
    *) echo "Unsupported OS: ${OS}"; exit 1;;
esac

# Detect Architecture
ARCH="$(uname -m)"
case "${ARCH}" in
    x86_64) ARCH='amd64';;
    arm64|aarch64) ARCH='arm64';;
    *) echo "Unsupported architecture: ${ARCH}"; exit 1;;
esac

echo "Detecting latest version..."
LATEST_TAG=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "${LATEST_TAG}" ]; then
    echo "Error: Could not find latest release for ${REPO}"
    exit 1
fi

echo "Latest version is ${LATEST_TAG}"

# Construct download URL
EXTENSION="tar.gz"
if [ "${OS}" = "windows" ]; then
    EXTENSION="zip"
fi

FILENAME="affectedpkgs-${LATEST_TAG}-${OS}-${ARCH}.${EXTENSION}"
URL="https://github.com/${REPO}/releases/download/${LATEST_TAG}/${FILENAME}"

# Create temporary directory
TMP_DIR=$(mktemp -d)
cd "${TMP_DIR}"

echo "Downloading ${URL}..."
curl -L -o "${FILENAME}" "${URL}"

echo "Extracting..."
if [ "${OS}" = "windows" ]; then
    unzip "${FILENAME}"
else
    tar -xzf "${FILENAME}"
fi

echo "Installing to ${INSTALL_DIR}..."
if [ "${OS}" = "windows" ]; then
    mv "${BINARY_NAME}.exe" "${INSTALL_DIR}/${BINARY_NAME}.exe"
else
    # Use sudo if we don't have write access to INSTALL_DIR
    if [ -w "${INSTALL_DIR}" ]; then
        mv "${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
        chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    else
        echo "Requesting sudo permissions to install to ${INSTALL_DIR}"
        sudo mv "${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
        sudo chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    fi
fi

echo "Cleaning up..."
rm -rf "${TMP_DIR}"

echo "Successfully installed ${BINARY_NAME} ${LATEST_TAG}!"
${BINARY_NAME} --help || echo "Installation complete. You may need to restart your terminal or add ${INSTALL_DIR} to your PATH."
