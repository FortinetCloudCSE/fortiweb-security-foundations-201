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
  value       = "https://${azurerm_public_ip.kalipip.ip_address}:8443"
  description = "Guacamole web UI. Sign in as guacadmin with the same password as the Terraform admin_password output; then open Lab Desktop (RDP to this VM)."
}

output "kali_native_rdp" {
  value       = "${azurerm_public_ip.kalipip.ip_address}:3389"
  description = "Native Microsoft Remote Desktop: connect here as labuser with admin_password (xrdp). Requires apply after NSG allows 3389."
}