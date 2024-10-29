
resource "azurerm_public_ip" "ubupip" {
  name                = "${var.username}_ubupip"
  location            = data.azurerm_resource_group.resourcegroup.location
  resource_group_name = data.azurerm_resource_group.resourcegroup.name
  allocation_method   = "Static"
}



resource "azurerm_network_interface" "ubunic" {
  name                = "${var.username}-node_ubunic"
  location            = data.azurerm_resource_group.resourcegroup.location
  resource_group_name = data.azurerm_resource_group.resourcegroup.name

  ip_configuration {
    name                          = "${var.username}_ipconfig"
    subnet_id                     = azurerm_subnet.protectedsubnet.id
    private_ip_address_allocation = "Static"
    private_ip_address            = "10.0.0.15"
    public_ip_address_id          = azurerm_public_ip.ubupip.id
  }
}

resource "azurerm_network_security_group" "ubu-nsg" {
  name                = "${var.username}-ubu_nsg"
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
    name                       = "allow-juice-inbound"
    priority                   = 102
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "3000"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  } 
}

resource "azurerm_network_interface_security_group_association" "ubunsg-association" {
  network_interface_id      = azurerm_network_interface.ubunic.id
  network_security_group_id = azurerm_network_security_group.ubu-nsg.id
}

resource "azurerm_linux_virtual_machine" "ubuntu" {
  name                  = "ubuntu-${var.username}"
  resource_group_name   = data.azurerm_resource_group.resourcegroup.name
  location              = data.azurerm_resource_group.resourcegroup.location
  size                  = var.vm_size
  admin_username        = var.admin_username
  admin_password        = var.admin_password
  disable_password_authentication = false
  custom_data = filebase64("${path.module}/ububoostrap.txt")

  network_interface_ids = [azurerm_network_interface.ubunic.id]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = var.ubupublisher
    offer     = var.ubuoffer
    sku       = var.ubusku
    version   = var.ubuvmversion
  }

}
