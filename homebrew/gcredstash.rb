require 'formula'

class Gcredstash < Formula
  VERSION = '0.2.7'

  homepage 'https://github.com/winebarrel/gcredstash'
  url "https://github.com/winebarrel/gcredstash/releases/download/v#{VERSION}/gcredstash-v#{VERSION}-darwin-amd64.gz"
  sha256 'c04c6de2c4464ce73e5c3c80aede33fd22a36df4c22b0b52c14d13add6a908f7'
  version VERSION
  head 'https://github.com/winebarrel/gcredstash.git', :branch => 'master'

  def install
    system "mv gcredstash-v#{VERSION}-darwin-amd64 gcredstash"
    bin.install 'gcredstash'
  end
end
