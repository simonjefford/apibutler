# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.require_version ">= 1.4"

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.box = "ratelimit "

  config.vm.network :forwarded_port, guest: 8080, host: 8081
  config.vm.network :forwarded_port, guest: 4000, host: 4001
  config.vm.network :forwarded_port, guest: 35729, host: 35729

  config.vm.provider :virtualbox do |v, override|
    override.vm.box_url = "http://sjjvagrantboxes.s3.amazonaws.com/ratelimit.box"
    v.customize ["modifyvm", :id, "--memory", "2048"]
  end
end
