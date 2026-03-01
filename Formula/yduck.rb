class Yduck < Formula
  desc "Mac dev environment & AI coding tools setup CLI"
  homepage "https://github.com/tc6-01/YangDuck"
  version "${VERSION}"

  on_macos do
    if Hardware::CPU.arm?
      url "https://gh-proxy.com/https://github.com/tc6-01/YangDuck/releases/download/v${VERSION}/yduck-darwin-arm64"
      sha256 "${SHA256_ARM64}"
    else
      url "https://gh-proxy.com/https://github.com/tc6-01/YangDuck/releases/download/v${VERSION}/yduck-darwin-amd64"
      sha256 "${SHA256_AMD64}"
    end
  end

  def install
    binary = Dir["yduck-darwin-*"].first
    bin.install binary => "yduck"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/yduck --version 2>&1", 0)
  end
end
