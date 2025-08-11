
data "azurerm_resource_group" "resourcegroup" {
  name = "${var.username}-http101-workshop"
}

resource "azurerm_virtual_network" "linuxvmnetwork" {
  name                = "${var.username}-http101_network"
  address_space       = ["10.0.0.0/24"]
  location            = data.azurerm_resource_group.resourcegroup.location
  resource_group_name = data.azurerm_resource_group.resourcegroup.name
}

resource "azurerm_subnet" "protectedsubnet" {
  name                 = "protected_subnet"
  resource_group_name  = data.azurerm_resource_group.resourcegroup.name
  virtual_network_name = azurerm_virtual_network.linuxvmnetwork.name
  address_prefixes     = ["10.0.0.0/26"]
}

resource "azurerm_public_ip" "kalipip" {
  name                = "${var.username}_kalipip"
  location            = data.azurerm_resource_group.resourcegroup.location
  resource_group_name = data.azurerm_resource_group.resourcegroup.name
  allocation_method   = "Static"
}

resource "azurerm_network_interface" "nic" {
  name                = "${var.username}-node_nic"
  location            = data.azurerm_resource_group.resourcegroup.location
  resource_group_name = data.azurerm_resource_group.resourcegroup.name

  ip_configuration {
    name                          = "${var.username}_ipconfig"
    subnet_id                     = azurerm_subnet.protectedsubnet.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.kalipip.id
  }
}

resource "azurerm_network_security_group" "kali-nsg" {
  name                = "${var.username}-kali_nsg"
  location            = data.azurerm_resource_group.resourcegroup.location
  resource_group_name = data.azurerm_resource_group.resourcegroup.name

  security_rule {
    name                       = "allow-ssh-inbound"
    priority                   = 101
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "22"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }

  security_rule {
    name                       = "allow-alt-https-inbound"
    priority                   = 102
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "8443"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
}

resource "azurerm_network_interface_security_group_association" "nsg-association" {
  network_interface_id      = azurerm_network_interface.nic.id
  network_security_group_id = azurerm_network_security_group.kali-nsg.id
}

resource "azurerm_linux_virtual_machine" "kalivm" {
  name                            = "kali-${var.username}"
  resource_group_name             = data.azurerm_resource_group.resourcegroup.name
  location                        = data.azurerm_resource_group.resourcegroup.location
  size                            = var.kali_vm_size
  admin_username                  = var.admin_username
  admin_password                  = var.admin_password
  disable_password_authentication = false

  # Use bootstrap.txt as cloud-init configuration
  custom_data = filebase64("bootstrap.txt")

  network_interface_ids = [azurerm_network_interface.nic.id]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Premium_LRS"
  }

  plan {
    publisher = var.publisher
    name      = var.sku
    product   = var.offer
  }

  source_image_reference {
    publisher = var.publisher
    offer     = var.offer
    sku       = var.sku
    version   = var.vmversion
  }

  depends_on = [
    azurerm_network_interface_security_group_association.nsg-association
  ]

  tags = {
    Name    = "kali-${var.username}"
    Purpose = "http101-workshop"
  }
}

output "kali_vm_public_ip-1" {
  value       = azurerm_public_ip.kalipip.ip_address
  description = "Public IP of the Kali Linux VM"
}

# Copy Go source files to the VM using file provisioner
resource "null_resource" "copy_go_files" {
  connection {
    type     = "ssh"
    host     = azurerm_linux_virtual_machine.kalivm.public_ip_address
    user     = var.admin_username
    password = var.admin_password
    timeout  = "5m"
  }

  # Create directory and set permissions first
  provisioner "remote-exec" {
    inline = [
      "echo '=== DEBUG: Starting directory creation ==='",
      "pwd",
      "whoami",
      "sudo mkdir -p /opt/mltool",
      "echo 'Directory created, checking permissions...'",
      "ls -la /opt/",
      "sudo chmod 755 /opt/mltool",
      "echo 'Permissions set, checking again...'",
      "ls -la /opt/",
      "sudo chown labuser:labuser /opt/mltool",
      "echo 'Ownership changed, final check...'",
      "ls -la /opt/",
      "echo '=== DEBUG: Directory setup complete ==='"
    ]
  }

  provisioner "file" {
    source      = "bots.go"
    destination = "/opt/mltool/bots.go"
  }

  provisioner "file" {
    source      = "ml2.go"
    destination = "/opt/mltool/ml2.go"
  }

  provisioner "file" {
    source      = "ml-mix.go"
    destination = "/opt/mltool/ml-mix.go"
  }

  # Verify files were copied and set permissions
  provisioner "remote-exec" {
    inline = [
      "echo '=== DEBUG: Verifying file copy ==='",
      "echo 'Current working directory:'",
      "pwd",
      "echo 'Checking if files exist:'",
      "ls -la /opt/mltool/",
      "echo 'File contents check (first 5 lines):'",
      "head -5 /opt/mltool/bots.go || echo 'bots.go not found or empty'",
      "head -5 /opt/mltool/ml2.go || echo 'ml2.go not found or empty'",
      "head -5 /opt/mltool/ml-mix.go || echo 'ml-mix.go not found or empty'",
      "echo 'Setting file permissions...'",
      "chmod 644 /opt/mltool/bots.go",
      "chmod 644 /opt/mltool/ml2.go",
      "chmod 644 /opt/mltool/ml-mix.go",
      "echo 'Final file check:'",
      "ls -la /opt/mltool/",
      "echo '=== DEBUG: File verification complete ==='"
    ]
  }

  # Trigger the build process (build script is now in bootstrap.txt)
  provisioner "remote-exec" {
    inline = [
      "echo '=== DEBUG: Triggering build process ==='",
      "echo 'Files copied successfully, triggering build...'",
      "echo 'Build script will run after Go installation completes'",
      "echo '=== DEBUG: Build trigger complete ==='"
    ]
  }

  depends_on = [azurerm_linux_virtual_machine.kalivm]
}