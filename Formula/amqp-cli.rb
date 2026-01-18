class AmqpCli < Formula
  desc "CLI tool for publishing and consuming RabbitMQ messages"
  homepage "https://github.com/zbum/amqp-cli"
  version "0.1.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/zbum/amqp-cli/releases/download/v#{version}/amqp-cli-darwin-arm64.tar.gz"
      sha256 "c11e1bf5547a1d05a116cb4f14ceb6e4c959f9131351f92c7ae723b14959cd24"
    else
      url "https://github.com/zbum/amqp-cli/releases/download/v#{version}/amqp-cli-darwin-amd64.tar.gz"
      sha256 "fffb049f5ee69319b340601027740b3d6c763e97b541b446c616666dbee16858"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/zbum/amqp-cli/releases/download/v#{version}/amqp-cli-linux-arm64.tar.gz"
      sha256 "aa6a8abae7ec6fee92202d71877922322d985bb11cc8239dc970776005d0fc75"
    else
      url "https://github.com/zbum/amqp-cli/releases/download/v#{version}/amqp-cli-linux-amd64.tar.gz"
      sha256 "036ee1ca869ff2927cd564ae4485fab58b5fc524f354c0ffa7ebdad089199ddb"
    end
  end

  def install
    bin.install "amqp-cli"
  end

  test do
    assert_match "amqp-cli", shell_output("#{bin}/amqp-cli --help")
  end
end
