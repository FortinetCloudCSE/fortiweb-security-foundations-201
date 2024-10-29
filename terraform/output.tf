output "kali_linux_PublicIP" {
  value       = azurerm_public_ip.kalipip.ip_address
}

output "ubuntu_PublicIP" {
  value       = azurerm_public_ip.ubupip.ip_address
}

output "ssh-username" {
  value       = azurerm_linux_virtual_machine.kalivm.admin_username 
}

output "webui-username" {
  value       = "guacadmin"
}

output "password" {
  value       = azurerm_linux_virtual_machine.kalivm.admin_password
  sensitive = true
}

output "login-url" {
  value = "https://${azurerm_public_ip.kalipip.ip_address}:8443"
}