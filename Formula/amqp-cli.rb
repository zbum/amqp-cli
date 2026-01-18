class AmqpCli < Formula
  desc "CLI tool for publishing and consuming RabbitMQ messages"
  homepage "https://github.com/zbum/amqp-cli"
  version "0.1.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/zbum/amqp-cli/releases/download/v#{version}/amqp-cli-darwin-arm64.tar.gz"
      sha256 "254f6446bd23c55a4bf90ae3a47e467adc437ddfe1c9dce4bc10d76026c777d1"
    else
      url "https://github.com/zbum/amqp-cli/releases/download/v#{version}/amqp-cli-darwin-amd64.tar.gz"
      sha256 "529d7a0e25b0835593fd00bc1f79e3ce54034807d26b04c1c0a1476558975fe2"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/zbum/amqp-cli/releases/download/v#{version}/amqp-cli-linux-arm64.tar.gz"
      sha256 "2288bfd8f9fcb349dd9546dbdc12be8830f3f8337535b9ec5c9df737b7372ee4"
    else
      url "https://github.com/zbum/amqp-cli/releases/download/v#{version}/amqp-cli-linux-amd64.tar.gz"
      sha256 "0d2820d308d910919a1f05998cd3abcd45a65bf2156d19c8e36cc2384a5885f7"
    end
  end

  def install
    bin.install "amqp-cli"
  end

  test do
    assert_match "amqp-cli", shell_output("#{bin}/amqp-cli --help")
  end
end
