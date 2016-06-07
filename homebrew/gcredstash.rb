require 'formula'

class Gcredstash < Formula
  VERSION = '0.1.1'

  homepage 'https://github.com/winebarrel/gcredstash'
  url "https://github.com/winebarrel/gcredstash/releases/download/v#{VERSION}/gcredstash-v#{VERSION}-darwin-amd64.gz"
  sha256 '227b51a38a50629c4fc695dbc2f36943f31603f33e6bdf3c20dfee8ad91c53d4'
  version VERSION
  head 'https://github.com/winebarrel/gcredstash.git', :branch => 'master'

  def install
    system "mv gcredstash-v#{VERSION}-darwin-amd64 gcredstash"
    bin.install 'gcredstash'
  end
end
