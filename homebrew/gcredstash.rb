require 'formula'

class Gcredstash < Formula
  VERSION = '0.3.5'

  homepage 'https://github.com/winebarrel/gcredstash'
  url "https://github.com/winebarrel/gcredstash/releases/download/v#{VERSION}/gcredstash-v#{VERSION}-darwin-amd64.gz"
  sha256 'bd3a8700e75ce0e0ae43da7dc5ec1170b7a1a683db1a50bd44527f10fc653129'
  version VERSION
  head 'https://github.com/winebarrel/gcredstash.git', :branch => 'master'

  def install
    system "mv gcredstash-v#{VERSION}-darwin-amd64 gcredstash"
    bin.install 'gcredstash'
  end
end
