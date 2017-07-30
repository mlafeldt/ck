class Cktool < Formula
  desc     "The ConvertKit Tool"
  homepage "https://github.com/mlafeldt/ck"
  version  "%VERSION%"
  url      "https://github.com/mlafeldt/ck/releases/download/v#{version}/ck_darwin_amd64"
  sha256   "%SHA%"

  bottle :unneeded

  def install
    bin.install "ck_darwin_amd64" => "ck"
  end

  test do
    system "#{bin}/ck --version"
  end
end
