#!/bin/bash

set -euxo pipefail

version=$(cat release/version)

darwin_cli_sha256=$(openssl dgst -sha256 release/ironbird-darwin | cut -d ' ' -f 2)
linux_cli_sha256=$(openssl dgst -sha256 release/ironbird-linux | cut -d ' ' -f 2)

pushd homebrew-tap
  cat <<EOF > ironbird.rb
class ControlTower < Formula
  desc "Test Concourse tasks using a YAML DSL for Ginkgo/Gomega"
  homepage "https://www.engineerbetter.com"
  license "Apache-2.0"
  version "${version}"

  if OS.mac?
    url "https://github.com/EngineerBetter/ironbird/releases/download/#{version}/ironbird-darwin"
    sha256 "${darwin_cli_sha256}"
  elsif OS.linux?
    url "https://github.com/EngineerBetter/ironbird/releases/download/#{version}/ironbird-linux"
    sha256 "${linux_cli_sha256}"
  end

  depends_on :arch => :x86_64

  def install
    binary_name = "ironbird"
    if OS.mac?
      bin.install "ironbird-darwin" => binary_name
    elsif OS.linux?
      bin.install "ironbird-linux" => binary_name
    end
  end

  test do
    system "#{bin}/ironbird --help"
  end
end
EOF
popd

