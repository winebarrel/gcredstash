require 'formula'

class Gcredstash < Formula
  VERSION = '0.2.4'

  homepage 'https://github.com/winebarrel/gcredstash'
  url "https://github.com/winebarrel/gcredstash/releases/download/v#{VERSION}/gcredstash-v#{VERSION}-darwin-amd64.gz"
  sha256 '4c80cf2a55ff17335f470ca603eecae8532fa75509e72ea412a3978ea3e33442'
  version VERSION
  head 'https://github.com/winebarrel/gcredstash.git', :branch => 'master'

  def install
    system "mv gcredstash-v#{VERSION}-darwin-amd64 gcredstash"
    bin.install 'gcredstash'
  end
end
