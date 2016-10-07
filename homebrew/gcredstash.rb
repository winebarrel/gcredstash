require 'formula'

class Gcredstash < Formula
  VERSION = '0.2.9'

  homepage 'https://github.com/winebarrel/gcredstash'
  url "https://github.com/winebarrel/gcredstash/releases/download/v#{VERSION}/gcredstash-v#{VERSION}-darwin-amd64.gz"
  sha256 '8d5fa90ac6b0cd0fc8c927ff11b7ce340439e1e82aad7bb7e08fd145f51edacb'
  version VERSION
  head 'https://github.com/winebarrel/gcredstash.git', :branch => 'master'

  def install
    system "mv gcredstash-v#{VERSION}-darwin-amd64 gcredstash"
    bin.install 'gcredstash'
  end
end
