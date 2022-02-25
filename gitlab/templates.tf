data "template_file" "init" {
  template = templatefile("${path.module}/templates/init.tftpl",
    {
      efs_id       = aws_efs_file_system.gitlab.id
      nfs_port     = "2049"
      nfs_dns_name = aws_efs_file_system.gitlab.dns_name
      sshd_port    = "8822"
    }
  )
}

data "template_file" "gen_inventory" {
  template = templatefile("${path.module}/templates/inventory.yml.tftpl",
    {
      name         = aws_instance.gitlab.tags_all["Name"]
      fqdn         = cloudflare_record.gitlab.hostname
      efs_dnsname  = aws_efs_file_system.gitlab.dns_name
      efs_port     = "2049"
      efs_id       = aws_efs_file_system.gitlab.id
      ansible_user = "ubuntu"
      ansible_port = "8822"
    }
  )
}

resource "local_file" "save_inventory" {
  content  = data.template_file.gen_inventory.rendered
  filename = "${path.module}/ansible/inventory/hosts.yml"
}
