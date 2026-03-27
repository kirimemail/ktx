#!/usr/bin/env bash
set -e

REPO="kirimemail/ktx"
INSTALL_DIR="/usr/local/bin"

get_os() {
    case "$(uname -s)" in
        Linux*)  echo "linux" ;;
        Darwin*) echo "darwin" ;;
        *)       echo " unsupported" && exit 1 ;;
    esac
}

get_arch() {
    case "$(uname -m)" in
        x86_64)  echo "amd64" ;;
        arm64)   echo "arm64" ;;
        aarch64) echo "arm64" ;;
        *)       echo " unsupported" && exit 1 ;;
    esac
}

get_ext() {
    if [ "$(get_os)" = "windows" ]; then
        echo ".exe"
    else
        echo ""
    fi
}

latest_tag() {
    curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/'
}

download() {
    local version="$1"
    local os="$2"
    local arch="$3"
    local ext="$4"
    local filename="ktx-${os}-${arch}${ext}"
    local url="https://github.com/${REPO}/releases/download/${version}/${filename}"
    
    echo "Downloading ${filename}..."
    curl -fsSL "$url" -o "${INSTALL_DIR}/ktx${ext}"
    chmod +x "${INSTALL_DIR}/ktx${ext}"
    
    if [ -f "${INSTALL_DIR}/ktx${ext}.sha256sum" ]; then
        echo "Verifying checksum..."
        cd "${INSTALL_DIR}"
        sha256sum -c "ktx${ext}.sha256sum"
        rm "ktx${ext}.sha256sum"
    fi
    
    echo "Installed ktx${ext} to ${INSTALL_DIR}"
}

checksum_url() {
    local version="$1"
    local os="$2"
    local arch="$3"
    local ext="$4"
    local filename="ktx-${os}-${arch}${ext}"
    local url="https://github.com/${REPO}/releases/download/${version}/${filename}.sha256sum"
    
    echo "Downloading checksums..."
    curl -fsSL "$url" -o "${INSTALL_DIR}/ktx${ext}.sha256sum"
}

install() {
    local version="${1:-$(latest_tag)}"
    local os=$(get_os)
    local arch=$(get_arch)
    local ext=$(get_ext)
    
    echo "Installing ktx ${version} for ${os}/${arch}..."
    
    download "$version" "$os" "$arch" "$ext"
    
    echo "Done!"
    ktx --version || true
}

uninstall() {
    rm -f "${INSTALL_DIR}/ktx"
    rm -f "${INSTALL_DIR}/ktx.exe"
    echo "Uninstalled ktx"
}

case "${1:-install}" in
    install)
        install "$2"
        ;;
    uninstall)
        uninstall
        ;;
    latest-tag)
        latest_tag
        ;;
    *)
        echo "Usage: $0 {install|uninstall|latest-tag}"
        exit 1
        ;;
esac
